-- CT Interactive / Nordic Song
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{10, 6, 7, 2, 11, 4, 8, 11, 4, 9, 5, 7, 6, 10, 8, 2, 9, 5, 11, 11, 11, 11, 10, 10, 10, 10, 2, 4, 4, 4, 4, 2, 5, 5, 5, 5, 8, 8, 8, 8, 9, 9, 9, 9, 3, 3, 3, 3, 7, 7, 7, 7, 6, 6, 6, 6, 3, 3, 3, 3},
	{11, 4, 9, 5, 7, 4, 10, 6, 7, 9, 8, 5, 11, 8, 6, 10, 1, 3, 3, 3, 3, 1, 1, 1, 1, 4, 4, 4, 4, 6, 6, 6, 6, 3, 3, 3, 3, 9, 9, 9, 9, 7, 7, 7, 7, 11, 11, 11, 11, 5, 5, 5, 5, 10, 10, 10, 10, 8, 8, 8, 8},
	{10, 1, 11, 5, 7, 4, 10, 5, 11, 4, 9, 2, 7, 6, 8, 9, 2, 8, 6, 6, 6, 6, 6, 9, 9, 9, 9, 5, 5, 5, 5, 7, 7, 7, 7, 10, 10, 10, 10, 3, 3, 3, 3, 2, 3, 3, 3, 3, 1, 1, 1, 1, 8, 8, 8, 8, 2, 11, 11, 11, 11, 4, 4, 4, 4},
	{11, 4, 10, 6, 11, 9, 5, 8, 6, 7, 4, 10, 9, 1, 7, 8, 5, 3, 3, 3, 3, 6, 6, 6, 6, 9, 9, 9, 9, 8, 8, 8, 8, 4, 4, 4, 4, 7, 7, 7, 7, 5, 5, 5, 5, 11, 11, 11, 11, 10, 10, 10, 10, 3, 3, 3, 3, 1, 1, 1, 1},
	{6, 8, 1, 9, 4, 7, 2, 10, 4, 11, 6, 9, 8, 5, 10, 2, 7, 5, 11, 1, 1, 1, 1, 2, 8, 8, 8, 8, 3, 3, 3, 3, 5, 5, 5, 5, 3, 3, 3, 3, 9, 9, 9, 9, 10, 10, 10, 10, 7, 7, 7, 7, 4, 4, 4, 4, 2, 11, 11, 11, 11, 6, 6, 6, 6},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{9, 4, 8, 5, 7, 11, 3, 10, 2, 9, 8, 11, 6, 7, 11, 4, 9, 2, 7, 10, 4, 7, 5, 11, 9, 5, 10, 6, 8, 4, 10, 8, 6, 3, 3, 3, 3},
	{7, 9, 10, 5, 11, 4, 8, 10, 5, 11, 6, 7, 1, 9, 10, 6, 9, 7, 8, 4, 10, 6, 11, 4, 8, 9, 7, 4, 11, 5, 8, 3, 3, 3, 3, 3},
	{9, 2, 11, 5, 8, 6, 10, 2, 9, 7, 4, 10, 11, 7, 4, 8, 5, 11, 4, 7, 9, 6, 7, 5, 10, 6, 9, 3, 11, 10, 8, 4, 11, 1, 8, 3, 3, 3, 3},
	{11, 3, 8, 6, 9, 11, 4, 10, 7, 5, 8, 7, 4, 9, 1, 10, 4, 8, 6, 7, 9, 5, 11, 4, 8, 10, 4, 11, 5, 9, 6, 10, 3, 3, 3, 3},
	{9, 4, 7, 5, 8, 11, 4, 8, 3, 7, 5, 10, 1, 11, 4, 10, 2, 11, 6, 7, 8, 6, 10, 9, 4, 11, 10, 2, 9, 7, 5, 9, 8, 6, 3, 3, 3, 3},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                     -- wild    (2, 3, 4, 5 reels only)
	[ 2] = {},                     -- scatter (1, 3, 5 reels only)
	[ 3] = {0, 10, 50, 200, 1000}, -- man
	[ 4] = {0, 0, 50, 150, 500},   -- woman
	[ 5] = {0, 0, 20, 100, 400},   -- owl
	[ 6] = {0, 0, 20, 100, 400},   -- dog
	[ 7] = {0, 0, 10, 50, 200},    -- ace
	[ 8] = {0, 0, 10, 50, 200},    -- king
	[ 9] = {0, 0, 5, 20, 100},     -- queen
	[10] = {0, 0, 5, 20, 100},     -- jack
	[11] = {0, 0, 5, 20, 100},     -- ten
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay, scat_fs = 5, 12 -- scatter pays and number of free spins awarded
local scat_min = 3 -- minimum scatters to win

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
		for symbol_id in pairs(PAYTABLE_LINE) do
			counts[symbol_id] = {}
			for i = 1, sx do counts[symbol_id][i] = 0 end
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

				-- 2-of-a-kind (XX---) EV
				local comb2 = c[1] * c[2] * (lens[3] - c[3]) * lens[4] * lens[5]
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
				current_comb * (lens[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

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
		print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%", sq, rtp_sym, rtp_fs))
	end
	local rtp_total
	do
		reels = reels_reg
		precalculate_reels()
		local rtp_line = calculate_line_ev() / reshuffles * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / reshuffles * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
