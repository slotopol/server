-- CT Interactive / Wild Clover
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{6, 6, 6, 10, 10, 10, 5, 5, 5, 4, 4, 4, 3, 3, 3, 8, 8, 8, 6, 6, 6, 3, 7, 7, 9, 9, 9, 2, 9, 9, 9, 10, 10, 10, 8, 8, 8, 2, 5, 5, 5, 4, 4, 7, 7, 7},
	{9, 9, 9, 2, 3, 3, 3, 4, 4, 4, 6, 6, 6, 7, 7, 7, 9, 9, 9, 5, 5, 10, 10, 10, 2, 8, 8, 8, 6, 6, 6, 8, 8, 8, 7, 7, 7, 1, 1, 1, 1, 5, 5, 5, 10, 10, 10, 4, 4, 3},
	{3, 3, 3, 10, 10, 10, 5, 5, 5, 7, 7, 2, 8, 8, 8, 9, 9, 9, 6, 6, 6, 8, 8, 8, 7, 7, 7, 1, 1, 1, 1, 4, 4, 3, 4, 4, 4, 10, 10, 10, 6, 6, 6, 2, 5, 5, 5, 9, 9, 9},
	{3, 10, 10, 10, 8, 8, 8, 5, 5, 5, 7, 7, 7, 2, 4, 4, 4, 9, 9, 9, 1, 1, 1, 1, 5, 5, 2, 8, 8, 8, 4, 4, 6, 6, 6, 7, 7, 7, 6, 6, 6, 10, 10, 10, 9, 9, 9, 3, 3, 3},
	{9, 9, 9, 5, 5, 5, 7, 7, 2, 9, 9, 9, 7, 7, 7, 5, 5, 5, 4, 4, 4, 3, 4, 4, 10, 10, 10, 8, 8, 8, 2, 3, 3, 3, 8, 8, 8, 6, 6, 6, 10, 10, 10, 6, 6, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                    -- wild (2, 3, 4 reels only)
	[ 2] = {},                    -- scatter
	[ 3] = {0, 4, 40, 100, 1000}, -- seven
	[ 4] = {0, 0, 30, 100, 300},  -- bell
	[ 5] = {0, 0, 10, 60, 200},   -- shoe
	[ 6] = {0, 0, 10, 60, 200},   -- coin
	[ 7] = {0, 0, 8, 35, 100},    -- peach
	[ 8] = {0, 0, 8, 35, 100},    -- apple
	[ 9] = {0, 0, 8, 35, 100},    -- plum
	[10] = {0, 0, 8, 35, 100},    -- cherry
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 3, 20, 500}
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
		local ev_sum = 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
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

		return ev_sum
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / N
	local rtp_scat = calculate_scat_ev() / N
	local rtp_total = rtp_line + rtp_scat
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
