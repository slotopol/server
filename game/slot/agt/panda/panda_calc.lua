-- AGT / Panda
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{6, 6, 6, 6, 7, 7, 7, 7, 7, 1, 3, 3, 3, 5, 4, 4, 2, 2, 2, 6, 6, 6, 6, 6, 5, 5, 5, 5, 8, 8, 8, 8, 8, 7, 7, 7, 7, 9, 3, 3, 3, 4, 4, 4, 4, 8, 8, 8, 8, 8, 5, 5, 5, 5},
	{7, 7, 7, 7, 7, 5, 5, 5, 5, 8, 8, 8, 8, 8, 6, 6, 6, 6, 8, 8, 8, 8, 8, 3, 3, 3, 5, 5, 5, 5, 9, 6, 6, 6, 6, 6, 4, 4, 4, 4, 2, 2, 2, 1, 4, 4, 7, 7, 7, 7, 5, 3, 3, 3},
	{1, 5, 5, 5, 5, 3, 3, 3, 4, 4, 4, 4, 9, 4, 4, 2, 2, 2, 6, 6, 6, 6, 6, 5, 6, 6, 6, 6, 8, 8, 8, 8, 8, 7, 7, 7, 7, 5, 5, 5, 5, 7, 7, 7, 7, 7, 3, 3, 3, 8, 8, 8, 8, 8},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = 500, -- wild
	[2] = 160, -- bonsai
	[3] = 80,  -- fish
	[4] = 40,  -- fan
	[5] = 20,  -- lamp
	[6] = 20,  -- pot
	[7] = 20,  -- flower
	[8] = 10,  -- button
	[9] = 0,   -- scatter
}

-- 3. CONFIGURATION
local sx, sy = 3, 3 -- screen width & height
local wild, scat = 1, 9 -- wild & scatter symbol IDs

-- Performs full RTP calculation for given reels
local function calculate(reels)
	assert(#reels == sx, "unexpected number of reels")
	-- Get number of total reshuffles and lengths of each reel.
	local reshuffles, lens = 1, {}
	for i, r in ipairs(reels) do
		reshuffles = reshuffles * #r
		lens[i] = #r
	end

	-- Count symbols occurrences on each reel
	local counts = {}
	for symbol_id in pairs(PAYTABLE_LINE) do
		counts[symbol_id] = {}
		for i = 1, sx do counts[symbol_id][i] = 0 end
	end
	for i, r in ipairs(reels) do
		for _, sym in ipairs(r) do
			counts[sym][i] = counts[sym][i] + 1
		end
	end

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev()
		local ev_sum = 0
		local w = counts[wild]
		local comb_w3 = w[1] * w[2] * w[3]

		-- Iterate through all symbols that pay on lines
		for symbol_id, pay in pairs(PAYTABLE_LINE) do
			if symbol_id ~= wild and symbol_id ~= scat then
				local s = counts[symbol_id]
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end
				local comb = c[1] * c[2] * c[3] - comb_w3
				ev_sum = ev_sum + comb * pay
			end
		end
		ev_sum = ev_sum + comb_w3 * PAYTABLE_LINE[wild]

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev()
		local c = counts[scat]
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= 1 then
					fs_sum = fs_sum + current_comb * scat_sum
					fs_num = fs_num + current_comb
				end
				return
			end
			-- Step 1: having a scatter on this reel
			find_scatter_combs(reel_index + 1, scat_sum + 1,
				current_comb * c[reel_index] * sy)
			-- Step 2: NOT having a scatter on this reel
			find_scatter_combs(reel_index + 1, scat_sum,
				current_comb * (lens[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

		return ev_sum, fs_sum, fs_num
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles * 100
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = sq * rtp_sym
	local rtp_total = rtp_sym + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
