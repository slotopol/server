-- CT Interactive / Halloween Fruits
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{5, 12, 4, 11, 12, 4, 7, 3, 8, 9, 3, 10, 5, 9, 10, 6, 11, 8, 6, 7, 12, 12, 12, 12, 10, 10, 10, 10, 3, 3, 3, 3, 3, 3, 9, 9, 9, 9, 8, 8, 8, 8, 5, 5, 5, 5, 6, 6, 6, 6, 4, 4, 4, 4, 11, 11, 11, 11, 7, 7, 7, 7},
	{10, 4, 9, 1, 2, 1, 1, 1, 11, 6, 12, 3, 9, 5, 10, 6, 8, 4, 7, 11, 5, 2, 8, 1, 1, 12, 2, 3, 7, 2, 2, 2, 1, 5, 5, 5, 5, 1, 1, 1, 1, 10, 10, 10, 10, 9, 9, 9, 9, 12, 12, 12, 12, 4, 4, 4, 4, 11, 11, 11, 11, 8, 8, 8, 8, 3, 3, 3, 3, 3, 3, 6, 6, 6, 6, 7, 7, 7, 7},
	{6, 9, 4, 8, 12, 5, 11, 3, 7, 1, 2, 1, 1, 1, 1, 1, 12, 5, 10, 6, 9, 2, 2, 2, 10, 4, 11, 3, 7, 1, 2, 8, 6, 6, 2, 6, 6, 11, 11, 11, 11, 5, 5, 5, 5, 8, 8, 8, 8, 1, 1, 1, 1, 7, 7, 7, 7, 12, 12, 12, 12, 4, 4, 4, 4, 9, 9, 9, 9, 3, 3, 3, 3, 3, 3, 10, 10, 10, 10},
	{7, 4, 9, 3, 11, 1, 2, 1, 9, 11, 6, 1, 1, 2, 2, 2, 1, 8, 5, 10, 1, 12, 6, 7, 2, 4, 8, 1, 2, 12, 5, 10, 3, 7, 7, 7, 7, 6, 6, 6, 6, 10, 10, 10, 10, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 11, 11, 11, 11, 1, 1, 1, 1, 12, 12, 12, 12, 5, 5, 5, 5, 8, 8, 8, 8, 9, 9, 9, 9},
	{5, 12, 6, 11, 5, 10, 6, 11, 1, 1, 9, 3, 7, 4, 9, 1, 7, 8, 1, 1, 1, 1, 12, 3, 10, 4, 8, 8, 8, 8, 8, 5, 5, 5, 5, 1, 1, 1, 1, 10, 10, 10, 10, 6, 6, 6, 6, 11, 11, 11, 11, 12, 12, 12, 12, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 9, 9, 9, 9, 7, 7, 7, 7},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{11, 8, 4, 9, 5, 11, 3, 10, 7, 4, 12, 5, 8, 3, 10, 6, 9, 7, 6, 12, 5, 5, 5, 5, 6, 6, 6, 6, 8, 8, 8, 8, 4, 4, 4, 4, 7, 7, 7, 7, 11, 11, 11, 11, 9, 9, 9, 9, 10, 10, 10, 10, 12, 12, 12, 12, 3, 3, 3, 3, 3, 3},
	{1, 2, 1, 8, 2, 12, 2, 3, 7, 5, 11, 4, 9, 1, 12, 6, 8, 7, 4, 10, 5, 9, 6, 10, 3, 11, 1, 2, 1, 2, 4, 4, 4, 4, 4, 4, 1, 1, 1, 1, 12, 12, 12, 12, 5, 5, 5, 5, 8, 8, 8, 8, 10, 10, 10, 10, 6, 6, 6, 6, 11, 11, 11, 11, 7, 7, 7, 7, 9, 9, 9, 9, 3, 3, 3, 3, 3, 3},
	{5, 12, 8, 5, 11, 8, 6, 10, 3, 12, 4, 7, 2, 9, 2, 1, 1, 1, 2, 2, 1, 2, 1, 1, 1, 11, 4, 9, 3, 10, 6, 7, 4, 4, 4, 4, 12, 12, 12, 12, 8, 8, 8, 8, 5, 5, 5, 5, 9, 9, 9, 9, 10, 10, 10, 10, 1, 1, 1, 1, 3, 3, 3, 3, 3, 3, 7, 7, 7, 7, 6, 6, 6, 6, 11, 11, 11, 11},
	{11, 4, 12, 3, 11, 5, 9, 2, 1, 9, 1, 2, 1, 7, 4, 8, 2, 1, 2, 12, 6, 10, 5, 7, 2, 1, 1, 1, 8, 3, 10, 6, 4, 4, 4, 4, 12, 12, 12, 12, 3, 3, 3, 3, 3, 3, 8, 8, 8, 8, 10, 10, 10, 10, 7, 7, 7, 7, 11, 11, 11, 11, 5, 5, 5, 5, 1, 1, 1, 1, 6, 6, 6, 6, 9, 9, 9, 9},
	{4, 11, 6, 9, 7, 10, 5, 11, 1, 1, 1, 8, 5, 12, 6, 9, 3, 7, 4, 10, 1, 1, 1, 8, 3, 12, 7, 7, 7, 7, 3, 3, 3, 3, 3, 3, 6, 6, 6, 6, 4, 4, 4, 4, 4, 1, 1, 1, 1, 8, 8, 8, 8, 12, 12, 12, 12, 11, 11, 11, 11, 5, 5, 5, 5, 9, 9, 9, 9, 10, 10, 10, 10},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                  -- wild (2, 3, 4, 5 reels only)
	[ 2] = {},                  -- scatter (2, 3, 4 reels only)
	[ 3] = {0, 0, 20, 50, 300}, -- witch
	[ 4] = {0, 0, 15, 30, 100}, -- cat
	[ 5] = {0, 0, 15, 30, 100}, -- banana
	[ 6] = {0, 0, 15, 30, 100}, -- grape
	[ 7] = {0, 0, 10, 15, 50},  -- apple
	[ 8] = {0, 0, 10, 15, 50},  -- melon
	[ 9] = {0, 0, 10, 15, 30},  -- orange
	[10] = {0, 0, 10, 15, 30},  -- lemon
	[11] = {0, 0, 10, 15, 30},  -- plum
	[12] = {0, 0, 10, 15, 30},  -- cherry
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 0, 3, 5}
local scat_fs = 15 -- number of free spins awarded on scatters wins
local scat_min = 4 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs

-- Performs full RTP calculation for given reels
local function calculate(reels_reg, reels_bon)
	assert(#reels_reg == sx, "unexpected number of regular reels")
	assert(#reels_bon == sx, "unexpected number of bonus reels")

	local reels
	local reshuffles, lens
	local counts

	-- Reels precalculations
	local function precalculate_reels()
		-- Get number of total reshuffles and lengths of each reel.
		reshuffles, lens = 1, {}
		for i, r in ipairs(reels) do
			reshuffles = reshuffles * #r
			lens[i] = #r
		end

		-- Count symbols occurrences on each reel
		counts = {}
		for sym_id in pairs(PAYTABLE_LINE) do
			counts[sym_id] = {}
			for i = 1, sx do counts[sym_id][i] = 0 end
		end
		for i, r in ipairs(reels) do
			for _, sym in ipairs(r) do
				counts[sym][i] = counts[sym][i] + 1
			end
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
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- 1. Preliminary calculation: Determine how many ways
		-- N scatter symbols can appear on each reel.
		-- ways_on_reel[reel_idx][scat_count] = number of combinations
		-- that yield scat_count scatters
		local ways_on_reel = {}
		for i, reel in ipairs(reels) do
			ways_on_reel[i] = {}
			for count = 0, sy do
				ways_on_reel[i][count] = 0
			end
			for stop_idx = 0, lens[i] - 1 do
				local count = 0
				for h = 0, sy - 1 do
					local symbol_idx = (stop_idx + h) % lens[i] + 1
					if reel[symbol_idx] == scat then
						count = count + 1
					end
				end
				ways_on_reel[i][count] = ways_on_reel[i][count] + 1
			end
		end

		-- 2. Recursive traversal of the combination tree to sum the EV
		-- reel_idx: current reel (1-5)
		-- scat_sum: the sum of scatters on the previous reels
		-- current_comb: the number of combinations that lead to this state
		local function find_scatter_ev_recursive(reel_idx, scat_sum, current_comb)

			-- Base case: all 5 reels have been processed
			if reel_idx > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[math.min(scat_sum, #PAYTABLE_SCAT)]
					fs_sum = fs_sum + current_comb * scat_fs
					fs_num = fs_num + current_comb
				end
				return
			end

			-- Recursive step: iterate through all possible outcomes
			-- of the current reel (0, 1, 2, or 3 scatters)
			for scat_count, ways in pairs(ways_on_reel[reel_idx]) do
				if ways > 0 then
					find_scatter_ev_recursive(
						reel_idx + 1,
						scat_sum + scat_count,
						current_comb * ways
					)
				end
			end
		end
		find_scatter_ev_recursive(1, 0, 1)

		return ev_sum, fs_sum, fs_num
	end

	-- Execute calculation
	local rtp_fs
	do
		reels = reels_bon
		precalculate_reels()
		local rtp_line = calculate_line_ev() / reshuffles * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / reshuffles * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_fs = sq * rtp_sym
		print(string.format("*bonus reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%", sq, rtp_sym, rtp_fs))
	end
	local rtp_total
	do
		reels = reels_reg
		precalculate_reels()
		local rtp_line = calculate_line_ev() / reshuffles * 100
		local ev_sum, fs_sum = calculate_scat_ev()
		local rtp_scat = ev_sum / reshuffles * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_sum*scat_fs))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)