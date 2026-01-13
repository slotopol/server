-- CT Interactive / American Gigolo
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{7, 5, 3, 10, 2, 5, 3, 8, 6, 9, 8, 4, 10, 7, 9, 4, 6, 2, 10, 9, 9, 9, 9, 7, 7, 7, 7, 10, 10, 10, 10, 6, 6, 6, 6, 3, 3, 3, 3, 8, 8, 8, 8, 5, 5, 5, 5, 4, 4, 4, 4},
	{6, 4, 9, 2, 5, 3, 10, 7, 2, 9, 3, 8, 6, 1, 8, 4, 5, 1, 7, 10, 4, 4, 4, 4, 5, 5, 5, 5, 10, 10, 10, 10, 3, 3, 3, 3, 8, 8, 8, 8, 9, 9, 9, 9, 6, 6, 6, 6, 7, 7, 7, 7},
	{9, 3, 10, 5, 2, 8, 7, 1, 8, 3, 6, 7, 4, 9, 5, 1, 10, 4, 6, 2, 5, 5, 5, 5, 8, 8, 8, 8, 4, 4, 4, 4, 9, 9, 9, 9, 3, 3, 3, 3, 7, 7, 7, 7, 10, 10, 10, 10, 6, 6, 6, 6},
	{5, 8, 4, 10, 2, 9, 7, 1, 5, 3, 8, 1, 6, 3, 7, 9, 4, 10, 2, 6, 10, 10, 10, 10, 9, 9, 9, 9, 6, 6, 6, 6, 7, 7, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 8, 8, 8, 8, 3, 3, 3, 3},
	{1, 7, 9, 10, 8, 1, 5, 3, 6, 2, 8, 4, 10, 3, 7, 2, 9, 6, 10, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7, 9, 9, 9, 9, 8, 8, 8, 8, 3, 3, 3, 3, 10, 10, 10, 10, 4, 4, 4, 4},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                   --  1 wild (2, 3, 4, 5 reels only)
	[ 2] = {},                   --  2 scatter
	[ 3] = {0, 0, 50, 150, 500}, --  3 blonde
	[ 4] = {0, 0, 15, 50, 100},  --  4 brunette
	[ 5] = {0, 0, 10, 50, 100},  --  5 cat
	[ 6] = {0, 0, 10, 50, 100},  --  6 dog
	[ 7] = {0, 0, 10, 20, 100},  --  7 ace
	[ 8] = {0, 0, 10, 20, 100},  --  8 king
	[ 9] = {0, 0, 10, 20, 100},  --  9 queen
	[10] = {0, 0, 10, 20, 100},  -- 10 jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 5, 10, 50}
local scat_fs = 12 -- number of free spins awarded
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs

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

		-- Iterate through all symbols that pay on lines
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			if symbol_id ~= wild and #pays > 0 then
				local s = counts[symbol_id]
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end

				-- 5-of-a-kind (XXXXX) EV
				local comb5 = c[1] * c[2] * c[3] * c[4] * c[5]
				ev_sum = ev_sum + comb5 * pays[5]

				-- 4-of-a-kind (XXXX-) EV
				local comb4 = c[1] * c[2] * c[3] * c[4] * (lens[5] - c[5])
				ev_sum = ev_sum + comb4 * pays[4]

				-- 3-of-a-kind (XXX--) EV
				local comb3 = c[1] * c[2] * c[3] * (lens[4] - c[4]) * lens[5]
				ev_sum = ev_sum + comb3 * pays[3]
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
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
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
