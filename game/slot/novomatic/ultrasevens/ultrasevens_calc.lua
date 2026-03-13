-- Novomatic / Ultra Sevens
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 2, 6, 4, 5, 7, 4, 5, 6, 7, 4, 6, 5, 7, 6, 4, 7, 3, 5, 6, 7, 4, 5, 7, 6, 4, 5, 7, 6, 5, 4},
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 2, 7, 5, 4, 6, 5, 7, 4, 6, 7, 4, 6, 7, 5, 4, 3, 7, 5, 6, 4, 5, 6, 7, 5, 6, 4, 7, 6, 4, 5},
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 4, 7, 6, 5, 3, 7, 6, 4, 5, 2, 7, 6, 4, 5, 7, 4, 6, 5, 4, 7, 6, 4, 7, 5, 6, 4, 5, 6, 7, 5},
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 5, 4, 7, 5, 4, 7, 5, 6, 7, 4, 6, 5, 4, 6, 7, 4, 6, 5, 7, 6, 4, 5, 7, 3, 6, 4, 5, 7, 2, 6},
	{7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 5, 6, 4, 7, 5, 6, 4, 5, 6, 2, 3, 7, 5, 6, 7, 4, 5, 7, 4, 5, 6, 7, 4, 6, 7, 5, 4, 6, 7, 4},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 10, 100, 1000, 10000}, -- 1 seven
	[2] = {0, 0, 40, 200, 500},      -- 2 melon
	[3] = {0, 0, 40, 200, 500},      -- 3 grapes
	[4] = {0, 0, 10, 50, 200},       -- 4 plum
	[5] = {0, 0, 10, 50, 200},       -- 5 orange
	[6] = {0, 0, 10, 50, 200},       -- 6 lemon
	[7] = {0, 5, 10, 50, 200},       -- 7 cherry
}

-- 3. CONFIGURATION
local sx, sy = 5, 4 -- grid width & height

-- 4. JACKPOT GROUPS
local GROUPS = {
	{1},
	{2, 3},
	{4, 5, 6, 7}
}

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

	-- Count windows with filled symbols
	local windows_counts = {}
	for x, reel in ipairs(reels) do
		local wc = {}
		local n = #reel
		local consecutive = 0
		local cur_sym = nil

		for i = 1, n + sy - 1 do
			local idx = (i - 1) % n + 1
			local sym = reel[idx]

			if sym == cur_sym then
				consecutive = consecutive + 1
				if consecutive >= sy then
					wc[sym] = (wc[sym] or 0) + 1
				end
			else
				cur_sym = sym
				consecutive = 1
			end
		end

		windows_counts[x] = wc
	end

	-- Count full screens
	local full_screens = {}
	for sym in pairs(PAYTABLE_LINE) do
		local combs = 1
		for i = 1, #reels do
			combs = combs * (windows_counts[i][sym] or 0)
		end
		full_screens[sym] = combs
	end

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev()
		local ev_sum = 0

		-- Iterate through all symbols that pay on lines
		for sym_id, pays in pairs(PAYTABLE_LINE) do
			local c = counts[sym_id]

			-- 5-of-a-kind (XXXXX) EV
			local comb5 = c[1] * c[2] * c[3] * c[4] * c[5]
			ev_sum = ev_sum + comb5 * pays[5]

			-- 4-of-a-kind (XXXX-) EV on left side
			local comb4 = c[1] * c[2] * c[3] * c[4] * (L[5] - c[5])
			ev_sum = ev_sum + comb4 * pays[4]

			-- 3-of-a-kind (XXX--) EV on left side
			local comb3 = c[1] * c[2] * c[3] * (L[4] - c[4]) * L[5]
			ev_sum = ev_sum + comb3 * pays[3]

			-- 2-of-a-kind (XX---) EV
			local comb2 = c[1] * c[2] * (L[3] - c[3]) * L[4] * L[5]
			ev_sum = ev_sum + comb2 * pays[2]
		end

		-- Subtract progressive jackpots
		for sym, combs in pairs(full_screens) do
			ev_sum = ev_sum - combs*PAYTABLE_LINE[sym][sx]
		end

		return ev_sum
	end

	local function calculate_jackpots()
		local results = {}
		for group_id, symbols in pairs(GROUPS) do
			local group_combs = 0
			for _, sym in ipairs(symbols) do
				group_combs = group_combs + full_screens[sym]
			end
			results[group_id] = group_combs
		end
		return results
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / N
	local rtp_scat = 0
	local rtp_total = rtp_line + rtp_scat
	local jp = calculate_jackpots()
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
	print(string.format("jackpots1: count %g, hit rate 1/%.12g", jp[1], N/jp[1]))
	print(string.format("jackpots2: count %g, hit rate 1/%.12g", jp[2], N/jp[2]))
	print(string.format("jackpots3: count %g, hit rate 1/%.12g", jp[3], N/jp[3]))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
