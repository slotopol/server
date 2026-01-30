-- AGT / Merry Christmas
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{5, 5, 5, 2, 1, 1, 1, 6, 6, 6, 6, 6, 4, 4, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 7, 2, 6, 6, 6, 6, 6, 6, 6, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 5, 5, 5, 3, 3, 3, 3, 2},
	{7, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 2, 7, 7, 7, 7, 7, 7, 7, 3, 3, 3, 3, 6, 6, 6, 6, 6, 6, 6, 7, 5, 5, 4, 4, 4, 4, 4, 4, 2, 1, 1, 1, 5, 5, 5, 5, 5, 5, 5},
	{4, 4, 4, 4, 4, 4, 2, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 6, 2, 6, 6, 6, 6, 6, 1, 1, 1, 3, 3, 3, 3, 5, 5, 5, 5, 5, 5, 5, 2, 7, 7, 7, 7, 7, 7, 7, 5, 5, 5},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = 500, -- snowman
	[2] = 0,   -- scatter
	[3] = 250, -- ice
	[4] = 100, -- sled
	[5] = 20,  -- house
	[6] = 10,  -- bell
	[7] = 5,   -- deer
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay, scat_fs = 10, 20 -- scatter pays and number of free spins awarded

-- 4. CONFIGURATION
local sx, sy = 3, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local scat_min = 3 -- minimum scatters to win

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
	local function calculate_line_ev_bon()
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

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev_reg()
		local ev_sum = 0

		-- Iterate through all symbols that pay on lines
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= scat then
				local s = counts[sym_id]
				local comb = s[1] * s[2] * s[3]
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
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * scat_pay
					fs_sum = fs_sum + current_comb * scat_fs
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
	local rtp_line_bon = calculate_line_ev_bon() / reshuffles * 100
	local rtp_line_reg = calculate_line_ev_reg() / reshuffles * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles * 100
	local rtp_sym_bon = rtp_line_bon + rtp_scat
	local rtp_sym_reg = rtp_line_reg + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = sq * rtp_sym_bon
	local rtp_total = rtp_sym_reg + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("*free games calculations*"))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line_bon, rtp_scat, rtp_sym_bon))
	print(string.format("*regular games calculations*"))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line_reg, rtp_scat, rtp_sym_reg))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym_reg, q, rtp_fs, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
