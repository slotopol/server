-- CT Interactive / Black Pharaoh
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{4, 9, 11, 6, 7, 2, 10, 7, 7, 7, 2, 11, 6, 3, 6, 9, 9, 9, 6, 2, 4, 9, 6, 11, 6, 5, 4, 5, 10, 2, 5, 8, 8, 8, 10, 5, 7, 5, 8, 7, 11, 11, 11, 10, 3, 9, 3, 8, 7, 3, 9, 10, 11, 10, 10, 10, 9, 8, 11, 8, 4, 5, 4},
	{5, 10, 8, 1, 11, 11, 11, 10, 6, 8, 10, 3, 6, 11, 4, 10, 3, 10, 5, 9, 5, 3, 6, 5, 3, 11, 9, 11, 8, 7, 7, 7, 8, 8, 8, 5, 11, 9, 9, 9, 11, 7, 9, 8, 9, 6, 1, 4, 5, 4, 1, 10, 10, 10, 4, 6, 7, 6, 4, 7, 8, 7, 9},
	{8, 8, 8, 7, 7, 7, 3, 7, 1, 10, 6, 3, 7, 4, 1, 4, 11, 4, 2, 7, 9, 6, 4, 9, 1, 8, 10, 5, 10, 10, 10, 6, 2, 10, 5, 9, 5, 2, 11, 8, 11, 11, 11, 2, 5, 9, 10, 6, 7, 4, 11, 5, 9, 9, 9, 3, 8, 9, 8, 10, 6, 11, 6, 5, 8, 11, 3},
	{6, 7, 7, 7, 3, 9, 9, 9, 7, 10, 5, 8, 5, 4, 6, 7, 5, 11, 8, 8, 8, 9, 8, 7, 3, 4, 5, 10, 10, 10, 5, 9, 10, 1, 8, 9, 7, 3, 11, 8, 6, 11, 4, 11, 1, 5, 10, 9, 4, 3, 6, 8, 6, 1, 11, 11, 11, 9, 6, 11, 10, 4, 10},
	{6, 7, 6, 11, 11, 11, 7, 5, 7, 5, 2, 5, 4, 7, 5, 10, 2, 8, 8, 8, 9, 4, 10, 4, 10, 11, 7, 7, 7, 9, 2, 6, 5, 8, 9, 2, 10, 6, 9, 9, 9, 8, 11, 1, 8, 5, 9, 3, 11, 4, 3, 11, 3, 10, 1, 8, 11, 6, 9, 4, 1, 10, 10, 10, 6, 3},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                     --  1 wild (2, 3, 4, 5 reels only)
	[ 2] = {},                     --  2 scatter (1, 3, 5 reels only)
	[ 3] = {0, 5, 100, 300, 1000}, --  3 cleopatra
	[ 4] = {0, 0, 50, 250, 500},   --  4 cat
	[ 5] = {0, 0, 30, 150, 400},   --  5 ankh
	[ 6] = {0, 0, 30, 150, 400},   --  6 eye
	[ 7] = {0, 0, 20, 100, 300},   --  7 ace
	[ 8] = {0, 0, 20, 100, 300},   --  8 king
	[ 9] = {0, 0, 10, 40, 200},    --  9 queen
	[10] = {0, 0, 10, 40, 200},    -- 10 jack
	[11] = {0, 0, 10, 40, 200},    -- 11 ten
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay = 6
local scat_fs = 10 -- number of free spins awarded
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs

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
				local comb4 = c[1] * c[2] * c[3] * c[4] * (L[5] - c[5])
				ev_sum = ev_sum + comb4 * pays[4]

				-- 3-of-a-kind (XXX--) EV
				local comb3 = c[1] * c[2] * c[3] * (L[4] - c[4]) * L[5]
				ev_sum = ev_sum + comb3 * pays[3]

				-- 2-of-a-kind (XX---) EV
				local comb2 = c[1] * c[2] * (L[3] - c[3]) * L[4] * L[5]
				ev_sum = ev_sum + comb2 * pays[2]
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
