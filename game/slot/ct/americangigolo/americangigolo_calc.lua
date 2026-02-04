-- CT Interactive / American Gigolo
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{8, 10, 10, 10, 10, 7, 7, 7, 7, 9, 5, 5, 5, 5, 3, 2, 8, 10, 8, 8, 8, 8, 4, 4, 4, 4, 9, 7, 9, 3, 6, 5, 6, 9, 10, 4, 6, 3, 7, 6, 9, 4, 10, 7, 8, 10, 3, 3, 3, 3, 6, 6, 6, 6, 4, 3, 2, 5, 4, 5, 7, 9, 9, 9, 9, 10, 8, 5, 8},
	{7, 7, 7, 7, 10, 1, 4, 4, 4, 4, 8, 10, 10, 10, 10, 5, 3, 1, 5, 8, 5, 5, 5, 5, 9, 3, 4, 9, 1, 9, 10, 7, 6, 3, 7, 4, 8, 8, 8, 8, 10, 5, 3, 3, 3, 3, 6, 6, 6, 6, 4, 5, 2, 9, 7, 8, 1, 8, 6, 10, 1, 9, 9, 9, 9, 8, 10, 2, 9, 3, 6, 1, 6, 7, 4},
	{3, 7, 4, 6, 1, 9, 5, 9, 1, 9, 10, 5, 7, 4, 1, 8, 6, 6, 6, 6, 8, 1, 8, 8, 8, 8, 9, 9, 9, 9, 1, 9, 10, 3, 4, 8, 5, 10, 6, 5, 4, 4, 4, 4, 6, 7, 1, 6, 5, 5, 5, 5, 10, 10, 10, 10, 7, 7, 7, 7, 3, 8, 7, 2, 4, 9, 3, 3, 3, 3, 10, 2, 3, 8, 10},
	{8, 10, 6, 1, 9, 6, 10, 6, 6, 6, 6, 10, 3, 8, 8, 8, 8, 5, 9, 2, 5, 6, 4, 1, 3, 4, 4, 4, 4, 8, 10, 10, 10, 10, 8, 6, 9, 9, 9, 9, 3, 8, 1, 8, 9, 7, 7, 7, 7, 4, 2, 9, 5, 7, 4, 7, 1, 3, 7, 3, 3, 3, 3, 1, 10, 5, 7, 10, 5, 5, 5, 5, 4, 9, 1},
	{10, 10, 10, 10, 3, 3, 3, 3, 4, 7, 1, 8, 4, 4, 4, 4, 6, 1, 10, 3, 9, 9, 9, 9, 6, 5, 4, 7, 3, 8, 9, 10, 9, 3, 5, 5, 5, 5, 2, 10, 4, 9, 6, 7, 2, 9, 7, 10, 7, 8, 1, 8, 5, 3, 6, 8, 1, 6, 6, 6, 6, 5, 1, 9, 8, 8, 8, 8, 5, 4, 7, 7, 7, 7, 10, 1},
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
local sx, sy = 5, 3 -- grid width & height
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

		-- Iterate through all symbols that pay on lines
		for sym_id, pays in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and #pays > 0 then
				local s = counts[sym_id]
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
	local rtp_line = calculate_line_ev() / reshuffles
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = sq * rtp_sym
	local rtp_total = rtp_sym + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym*100, q, rtp_fs*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
