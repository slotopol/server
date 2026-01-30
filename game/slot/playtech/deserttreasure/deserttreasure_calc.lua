-- Playtech / Desert Treasure
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{8, 2, 4, 5, 7, 10, 3, 9, 10, 6, 8, 4, 6, 9, 3, 7, 1, 5, 10, 9, 7, 4, 8, 6, 5, 3},
	{2, 10, 4, 5, 8, 7, 9, 3, 5, 1, 10, 3, 7, 4, 9, 8, 6, 4, 10, 5, 7, 6, 3, 9, 8, 6},
	{3, 7, 1, 3, 4, 5, 6, 3, 9, 10, 6, 2, 10, 5, 4, 7, 8, 4, 9, 8, 5, 9, 6, 10, 7, 8},
	{9, 4, 10, 7, 3, 4, 10, 6, 2, 8, 5, 3, 7, 6, 8, 1, 9, 4, 6, 8, 3, 5, 9, 7, 10, 5},
	{3, 9, 10, 6, 9, 3, 6, 9, 10, 7, 5, 8, 7, 1, 4, 7, 6, 2, 3, 5, 8, 4, 5, 8, 4, 10},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{9, 10, 4, 7, 1, 3, 7, 10, 6, 4, 5, 8, 9, 3, 6, 10, 2, 8, 7, 10, 8, 7, 5, 9, 3, 5, 9, 4, 8, 6, 5},
	{6, 7, 8, 6, 5, 1, 10, 9, 4, 10, 8, 9, 7, 8, 10, 3, 5, 10, 9, 8, 4, 5, 6, 4, 7, 5, 3, 7, 2, 3, 9},
	{7, 3, 5, 9, 10, 8, 4, 7, 9, 10, 2, 7, 6, 9, 7, 4, 3, 10, 8, 6, 1, 5, 10, 6, 5, 8, 3, 9, 8, 5, 4},
	{10, 9, 8, 6, 7, 4, 1, 8, 7, 3, 8, 10, 9, 7, 8, 6, 9, 4, 5, 10, 7, 5, 3, 2, 10, 5, 6, 9, 3, 5, 4},
	{8, 10, 5, 7, 6, 10, 5, 3, 4, 8, 9, 7, 5, 3, 9, 4, 1, 7, 9, 8, 7, 6, 8, 5, 9, 3, 10, 2, 4, 10, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 8, 80, 800, 8000}, -- wild
	[ 2] = {},                    -- scatter
	[ 3] = {0, 5, 40, 250, 1000}, -- shield
	[ 4] = {0, 0, 20, 75, 500},   -- swords
	[ 5] = {0, 0, 0, 50, 250},    -- lamp
	[ 6] = {0, 2, 10, 30, 150},   -- ligature1
	[ 7] = {0, 2, 10, 30, 150},   -- ligature2
	[ 8] = {0, 0, 5, 15, 75},     -- ligature3
	[ 9] = {0, 0, 5, 15, 75},     -- ligature4
	[10] = {0, 0, 0, 10, 50},     -- ligature5
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 4, 40, 400}

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 10, 25, 50}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 3 -- minimum scatters to win
local mfs = 3 -- multiplier on free spins

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
		local wpays = PAYTABLE_LINE[wild]

		for sym_id, pays in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and #pays > 0 then
				local s = counts[sym_id]
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end

				for n = line_min, sx do
					local payout = pays[n]
					if payout > 0 then
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

						local better_wilds = 0
						local wn_min = nil
						for wn = line_min, n do
							if wpays[wn] >= payout then
								wn_min = wn
								break
							end
						end
						if wn_min then
							local bw = 1
							for i = 1, sx do
								if i <= wn_min then
									bw = bw * w[i]
								elseif i <= n then
									bw = bw * c[i]
								elseif i == n + 1 then
									bw = bw * (lens[i] - c[i])
								else
									bw = bw * lens[i]
								end
							end
							better_wilds = bw
						end

						ev_sum = ev_sum + (combs_total - better_wilds) * payout
					end
				end
			end
		end

		-- Calculating wilds as a separate symbol
		for n = line_min, sx do
			local payout = wpays[n]
			if payout > 0 then
				-- 1. Count all "clean heads" from wilds of length exactly n
				local wc = 1
				for i = 1, sx do
					if i <= n then wc = wc * w[i]
					elseif i == n + 1 then wc = wc * (lens[i] - w[i])
					else wc = wc * lens[i] end
				end

				-- 2. Subtract the cases where this line of wilds is intercepted by the S symbol.
				local losses = 0
				if n < sx then
					for sym_id, pays in pairs(PAYTABLE_LINE) do
						if sym_id ~= wild and #pays > 0 then
							local s = counts[sym_id]
							local c = {}
							for i = 1, sx do c[i] = s[i] + w[i] end

							for sn = n + 1, sx do
								if pays[sn] > payout then
									local loss = 1
									for i = 1, sx do
										if i <= n then
											loss = loss * w[i]
										elseif i == n + 1 then
											loss = loss * s[i]
										elseif i <= sn then
											loss = loss * c[i]
										elseif i == sn + 1 then
											loss = loss * (lens[i] - c[i])
										else
											loss = loss * lens[i]
										end
									end
									losses = losses + loss
								end
							end
						end
					end
				end
				ev_sum = ev_sum + (wc - losses) * payout
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
		rtp_fs = mfs * sq * rtp_sym
		print(string.format("*bonus reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %g*sq*rtp(sym) = %g*%.5g*%.5g = %.6f%%", mfs, mfs, sq, rtp_sym, rtp_fs))
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
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
