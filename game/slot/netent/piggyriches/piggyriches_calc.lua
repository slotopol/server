-- NetEnt / Piggy Riches
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{4, 12, 11, 10, 6, 9, 8, 12, 11, 5, 12, 9, 7, 3, 8, 2, 12, 6, 10, 9, 11, 8, 4, 10, 9, 11, 12, 10, 11, 6, 12, 1, 8, 9, 7, 10, 11, 5, 7},
	{8, 9, 10, 12, 11, 3, 9, 12, 7, 8, 9, 12, 7, 11, 6, 12, 11, 7, 5, 4, 10, 9, 12, 10, 5, 8, 10, 12, 11, 4, 10, 2, 6, 9, 11, 8, 6, 11, 1},
	{4, 10, 1, 7, 9, 10, 11, 8, 10, 12, 8, 11, 12, 6, 11, 12, 7, 9, 5, 12, 11, 9, 12, 11, 2, 6, 12, 9, 10, 7, 11, 6, 10, 5, 8, 3, 4, 9, 8},
	{12, 7, 11, 8, 12, 3, 10, 9, 5, 7, 9, 4, 8, 11, 7, 10, 12, 6, 8, 11, 1, 6, 12, 2, 11, 10, 6, 5, 9, 4, 10, 9, 12, 10, 11, 9, 12, 11, 8},
	{8, 12, 6, 11, 10, 4, 12, 10, 11, 9, 5, 6, 10, 7, 9, 11, 7, 3, 2, 8, 11, 12, 4, 8, 9, 11, 10, 1, 9, 12, 8, 11, 7, 12, 6, 9, 10, 12, 5},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                    -- wild (on all reels)
	[ 2] = {},                    -- scatter
	[ 3] = {0, 5, 25, 300, 2000}, -- money bag
	[ 4] = {0, 0, 25, 150, 1000}, -- banknotes
	[ 5] = {0, 0, 20, 125, 750},  -- keys
	[ 6] = {0, 0, 20, 75, 400},   -- wallet
	[ 7] = {0, 0, 15, 75, 200},   -- piggy bank
	[ 8] = {0, 0, 15, 50, 125},   -- ace
	[ 9] = {0, 0, 10, 25, 100},   -- king
	[10] = {0, 0, 5, 20, 75},     -- queen
	[11] = {0, 0, 5, 15, 60},     -- jack
	[12] = {0, 0, 5, 10, 50},     -- ten
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 2, 4, 15, 100}
local scat_min = 1 -- minimum scatters to win

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 15, 15, 15}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local mfs = 3 -- multiplier on free spins

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

					-- Total combinations where combination length is EXACTLY n
					-- Formula: (C1 * C2 * ... * Cn) * (Lens[n+1] - Cn+1) * ...
					local combs_total = 1
					for i = 1, 5 do
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
					for i = 1, 5 do
						if i <= n then
							combs_no_wild = combs_no_wild * s[i]
						elseif i == n + 1 then
							combs_no_wild = combs_no_wild * (lens[i] - c[i])
						else
							combs_no_wild = combs_no_wild * lens[i]
						end
					end

					local combs_only_wild = 1
					for i = 1, 5 do
						if i <= n then
							combs_only_wild = combs_only_wild * w[i]
						elseif i == n + 1 then
							combs_only_wild = combs_only_wild * (lens[i] - c[i])
						else
							combs_only_wild = combs_only_wild * lens[i]
						end
					end

					local combs_with_wild = combs_total - combs_no_wild - combs_only_wild
					return (combs_no_wild * payout) + (combs_with_wild * payout * 3)
				end

				for n = 2, 5 do
					ev_sum = ev_sum + get_comb_ev(n, pays[n])
				end
			end
		end

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev(free_spins)
		local c = counts[scat]
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					if free_spins then
						fs_sum = fs_sum + current_comb * scat_sum
						fs_num = fs_num + current_comb
					else
						ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
						if FREESPIN_SCAT[scat_sum] > 0 then
							fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
							fs_num = fs_num + current_comb
						end
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
	local rtp_fs
	local rtp_line = calculate_line_ev() / reshuffles
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	do
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(true)
		local rtp_scat = ev_sum / reshuffles
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_fs = mfs * sq * rtp_sym
		print(string.format("*free games calculations*"))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", 5*rtp_line*100, 5*rtp_scat*100, 5*rtp_sym*100))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %g*sq*rtp(sym) = %g*%.5g*%.5g = %.6f%%", mfs, mfs, sq, rtp_sym*100, rtp_fs*100))
	end
	local rtp_total
	do
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(false)
		local rtp_scat = ev_sum / reshuffles
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular games calculations*"))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym*100, q, rtp_fs*100, rtp_total*100))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
