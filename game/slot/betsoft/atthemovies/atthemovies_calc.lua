-- BetSoft / At the Movies
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{7, 8, 5, 6, 7, 5, 3, 8, 6, 5, 8, 7, 5, 1, 7, 4, 8, 6, 5, 7, 4, 6, 10, 8, 7, 4, 6, 3, 8, 7, 2, 6, 8, 4},
	{5, 2, 7, 4, 6, 3, 8, 7, 5, 6, 10, 8, 3, 6, 8, 5, 7, 6, 8, 7, 6, 8, 7, 4, 6, 8, 5, 4, 6, 1, 7, 9, 8, 7},
	{8, 6, 3, 5, 8, 4, 6, 1, 7, 6, 8, 7, 5, 8, 7, 6, 8, 7, 4, 8, 5, 2, 7, 6, 4, 7, 5, 6, 3, 5, 8, 7, 4, 10},
	{7, 6, 8, 5, 6, 7, 4, 8, 7, 6, 5, 7, 3, 8, 4, 6, 8, 7, 5, 8, 2, 5, 7, 8, 6, 1, 8, 6, 7, 9, 3, 6, 10, 4},
	{5, 8, 6, 1, 5, 7, 8, 6, 10, 5, 7, 6, 5, 7, 6, 4, 7, 5, 8, 4, 7, 8, 4, 7, 3, 8, 6, 7, 3, 8, 2, 6, 8, 4},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 20, 200, 500, 1000}, -- oscar
	[ 2] = {0, 10, 100, 250, 500},  -- popcorn
	[ 3] = {0, 5, 50, 100, 200},    -- poster
	[ 4] = {0, 2, 25, 50, 100},     -- a
	[ 5] = {0, 0, 20, 40, 80},      -- dummy
	[ 6] = {0, 0, 15, 30, 60},      -- maw
	[ 7] = {0, 0, 10, 20, 40},      -- starship
	[ 8] = {0, 0, 5, 10, 20},       -- heart
	[ 9] = {},                      -- masks (2, 4 reels only)
	[10] = {},                      -- projector
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 2, 0, 0, 0}

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 8, 12, 20}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 9, 10 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 2 -- minimum scatters to win
local mw = 2 -- multiplier on wilds
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

		for sym_id, pays in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and #pays > 0 then
				local s = counts[sym_id]
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end

				-- Function to calculate combinations with multiplier logic
				local function get_comb_ev(n, payout)
					if payout <= 0 then return 0 end

					-- Total combinations
					local combs_total = 1
					for i = 1, sx do
						if i <= n then
							combs_total = combs_total * c[i]
						elseif i == n + 1 then
							combs_total = combs_total * (lens[i] - c[i])
						else
							combs_total = combs_total * lens[i]
						end
					end

					-- Combinations WITHOUT any wilds on reels
					local combs_no_wild = 1
					for i = 1, sx do
						if i <= n then
							combs_no_wild = combs_no_wild * s[i]
						elseif i == n + 1 then
							combs_no_wild = combs_no_wild * (lens[i] - c[i])
						else
							combs_no_wild = combs_no_wild * lens[i]
						end
					end

					local combs_with_wild = combs_total - combs_no_wild
					return (combs_no_wild + combs_with_wild * mw) * payout
				end

				for n = line_min, sx do
					ev_sum = ev_sum + get_comb_ev(n, pays[n])
				end
			end
		end

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev(free_spins)
		local c = counts[scat]
		local mm = free_spins and mfs or 1
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum] * mm
					if FREESPIN_SCAT[scat_sum] > 0 then
						fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
						fs_num = fs_num + current_comb
					end
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
	local rtp_fs = mfs * sq * rtp_sym
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
