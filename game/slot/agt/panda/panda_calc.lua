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
local sx, sy = 3, 3 -- grid width & height
local wild, scat = 1, 9 -- wild & scatter symbol IDs

-- Performs full RTP calculation for given reels
local function calculate(reels)
	assert(#reels == sx, "unexpected number of reels")

	-- Get number of total reshuffles and lengths of each reel.
	local N, L = 1, {}
	for i, r in ipairs(reels) do
		N = N * #r
		L[i] = #r
	end

	-- Count symbols occurrences on each reel
	local counts = {}
	for sym_id in pairs(PAYTABLE_LINE) do
		counts[sym_id] = {}
		for i = 1, sx do counts[sym_id][i] = 0 end
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
		ev_sum = ev_sum + comb_w3 * PAYTABLE_LINE[wild]

		-- Iterate through all symbols that pay on lines
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and sym_id ~= scat then
				local s = counts[sym_id]
				local comb = (s[1] + w[1]) * (s[2] + w[2]) * (s[3] + w[3]) - comb_w3
				ev_sum = ev_sum + comb * pay
			end
		end

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
				current_comb * (L[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

		return ev_sum, fs_sum, fs_num
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / N
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / N
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / N
	local sq = 1 / (1 - q)
	local rtp_fs = sq * rtp_sym
	local rtp_total = rtp_sym + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games hit rate: 1/%.5g", N/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym*100, q, rtp_fs*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
