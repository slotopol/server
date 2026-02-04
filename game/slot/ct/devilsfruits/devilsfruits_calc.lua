-- CT Interactive / Devil's Fruits
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{2, 6, 10, 0, 7, 3, 5, 0, 8, 4, 3, 0, 5, 1, 2, 0, 6, 9, 10, 0, 4, 5, 10, 0, 6, 9, 5, 0, 6, 3, 9, 0, 2, 8, 7, 0, 6, 8, 7, 0, 10, 8, 9, 0, 1, 4, 10, 0, 7, 2, 4, 0, 3, 9, 8, 0, 7, 0},
	{4, 10, 7, 0, 8, 1, 7, 0, 4, 3, 9, 0, 2, 8, 4, 0, 10, 8, 2, 0, 6, 10, 9, 0, 3, 5, 7, 0, 9, 5, 6, 0, 10, 5, 2, 0, 7, 6, 2, 0, 9, 10, 7, 0, 3, 5, 4, 0, 9, 3, 8, 0, 1, 6, 8, 0},
	{7, 6, 9, 0, 10, 3, 8, 0, 10, 9, 5, 0, 8, 2, 4, 0, 3, 5, 2, 0, 8, 6, 1, 0, 10, 7, 9, 0, 4, 6, 5, 0, 2, 10, 4, 0, 6, 7, 2, 0, 6, 10, 7, 0, 3, 5, 9, 0, 1, 7, 9, 0, 8, 4, 3, 0, 8, 0},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 0] = 0,    -- space
	[ 1] = 1500, -- wild
	[ 2] = 100,  -- seven
	[ 3] = 35,   -- pike
	[ 4] = 25,   -- bell
	[ 5] = 25,   -- orange
	[ 6] = 25,   -- plum
	[ 7] = 25,   -- bar3
	[ 8] = 20,   -- bar2
	[ 9] = 15,   -- bar1
	[10] = 10,   -- cherry
}

-- 3. CONFIGURATION
local sx = 3 -- grid width
-- Symbols names
local space = 0
local wild  = 1
local bar3  = 7
local bar2  = 8
local bar1  = 9

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
		local b1 = counts[bar1]
		local b2 = counts[bar2]
		local b3 = counts[bar3]

		-- Iterate through all symbols that pay on lines
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= space then
				local s = counts[sym_id]
				local comb = s[1] * s[2] * s[3]
				ev_sum = ev_sum + comb * pay
			end
		end
		-- 1 wild and any 2 symbols
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and sym_id ~= space then
				local s = counts[sym_id]
				local comb =
					w[1] * s[2] * s[3] +
					s[1] * w[2] * s[3] +
					s[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 2
			end
		end
		-- 2 wilds and any 1 symbol
		for sym_id, pay in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and sym_id ~= space then
				local s = counts[sym_id]
				local comb =
					w[1] * w[2] * s[3] +
					s[1] * w[2] * w[3] +
					w[1] * s[2] * w[3]
				ev_sum = ev_sum + comb * pay * 4
			end
		end
		-- any bars with wild
		do
			local comb =
				w[1] * b1[2] * b2[3] +
				w[1] * b1[2] * b3[3] +
				w[1] * b2[2] * b1[3] +
				w[1] * b2[2] * b3[3] +
				w[1] * b3[2] * b1[3] +
				w[1] * b3[2] * b2[3] +
				b1[1] * w[2] * b2[3] +
				b1[1] * w[2] * b3[3] +
				b2[1] * w[2] * b1[3] +
				b2[1] * w[2] * b3[3] +
				b3[1] * w[2] * b1[3] +
				b3[1] * w[2] * b2[3] +
				b1[1] * b2[2] * w[3] +
				b1[1] * b3[2] * w[3] +
				b2[1] * b1[2] * w[3] +
				b2[1] * b3[2] * w[3] +
				b3[1] * b1[2] * w[3] +
				b3[1] * b2[2] * w[3]
			ev_sum = ev_sum + comb * 5 * 2
		end
		-- any bar without wild
		do
			local b = {}
			for i = 1, sx do b[i] = b1[i] + b2[i] + b3[i] end
			local comb = b[1] * b[2] * b[3] -
				b1[1] * b1[2] * b1[3] -
				b2[1] * b2[2] * b2[3] -
				b3[1] * b3[2] * b3[3]
			ev_sum = ev_sum + comb * 5
		end

		return ev_sum
	end

	-- Execute calculation
	local rtp_total = calculate_line_ev() / N
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
	print(string.format("RTP = %.6f%%", rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
