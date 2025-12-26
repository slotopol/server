-- CT Interactive / Sapphire Lagoon
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{3, 7, 3, 9, 3, 8, 5, 7, 9, 4, 8, 6, 10, 4, 4, 6, 8, 10, 2, 6, 10, 5, 7, 9, 5, 9, 9, 9, 9, 4, 4, 4, 4, 7, 7, 7, 7, 8, 8, 8, 8, 5, 5, 5, 5, 3, 3, 3, 3, 6, 6, 6, 6, 10, 10, 10, 10, 2, 2, 2, 2},
	{4, 3, 2, 3, 8, 10, 6, 4, 5, 9, 10, 8, 4, 6, 7, 1, 9, 5, 6, 1, 8, 7, 9, 5, 10, 7, 10, 10, 10, 10, 2, 2, 2, 8, 8, 8, 8, 5, 5, 5, 5, 1, 1, 1, 1, 3, 3, 3, 3, 7, 7, 7, 7, 6, 6, 6, 6, 9, 9, 9, 9, 4, 4, 4, 4},
	{10, 10, 6, 1, 3, 4, 4, 5, 10, 8, 9, 5, 4, 9, 1, 7, 8, 3, 9, 5, 2, 7, 6, 6, 8, 7, 9, 9, 9, 9, 10, 10, 10, 10, 6, 6, 6, 6, 4, 4, 4, 4, 8, 8, 8, 8, 3, 3, 3, 3, 2, 2, 2, 5, 5, 5, 5, 1, 1, 1, 1, 7, 7, 7, 7},
	{6, 1, 7, 10, 9, 5, 8, 1, 5, 2, 10, 10, 4, 8, 3, 9, 4, 10, 9, 6, 7, 5, 7, 3, 6, 8, 4, 8, 8, 8, 8, 5, 5, 5, 5, 1, 1, 1, 1, 9, 9, 9, 9, 6, 6, 6, 6, 2, 2, 2, 4, 4, 4, 4, 7, 7, 7, 7, 3, 3, 3, 3, 10, 10, 10, 10},
	{1, 10, 7, 6, 5, 3, 1, 6, 8, 6, 4, 4, 2, 9, 10, 7, 9, 10, 3, 8, 9, 5, 8, 10, 5, 7, 4, 4, 4, 4, 4, 9, 9, 9, 9, 10, 10, 10, 10, 8, 8, 8, 8, 6, 6, 6, 6, 3, 3, 3, 3, 7, 7, 7, 7, 1, 1, 1, 1, 5, 5, 5, 5, 2, 2, 2},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 0, 0, 0},       -- wild (2, 3, 4, 5 reels only)
	[ 2] = {0, 0, 0, 0, 0},       -- scatter
	[ 3] = {0, 0, 20, 200, 1000}, -- man
	[ 4] = {0, 0, 15, 75, 150},   -- woman
	[ 5] = {0, 0, 5, 50, 150},    -- flask
	[ 6] = {0, 0, 5, 50, 150},    -- hook
	[ 7] = {0, 0, 5, 15, 100},    -- ace
	[ 8] = {0, 0, 5, 15, 100},    -- king
	[ 9] = {0, 0, 5, 15, 100},    -- queen
	[10] = {0, 0, 5, 15, 100},    -- jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 0, 0, 0, 3, 5, 10, 15, 20, 25, 30, 40, 50, 100}
local scat_min = 6 -- minimum scatters to win

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 0, 0, 0, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local mfs = 2 -- multiplier on free spins

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
			if symbol_id ~= wild and symbol_id ~= scat then
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
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
					fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
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
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles * 100
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = mfs * sq * rtp_sym
	local rtp_total = rtp_sym + q*rtp_fs
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
