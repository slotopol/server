-- CT Interactive / Treasure Kingdom
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{5, 9, 13, 4, 7, 10, 9, 8, 1, 7, 4, 6, 8, 13, 10, 5, 13, 11, 7, 2, 13, 5, 6, 11, 12, 10, 9, 12, 11, 10, 12, 3, 6, 1, 8, 4, 11, 2, 9, 3, 13, 12, 8},
	{2, 12, 4, 10, 7, 4, 11, 1, 13, 12, 2, 8, 3, 7, 13, 6, 9, 5, 10, 13, 12, 3, 11, 4, 8, 9, 1, 10, 8, 13, 9, 5, 12, 11, 9, 6, 5, 10, 13, 7, 8, 6, 11},
	{13, 10, 11, 6, 3, 10, 13, 8, 1, 12, 10, 5, 9, 7, 2, 9, 4, 8, 10, 4, 11, 8, 5, 11, 3, 12, 1, 9, 6, 4, 8, 9, 13, 5, 12, 13, 7, 2, 13, 6, 12, 7, 11},
	{8, 9, 11, 10, 7, 4, 12, 3, 13, 4, 8, 12, 1, 10, 13, 5, 8, 9, 10, 2, 11, 12, 10, 2, 9, 13, 7, 5, 13, 9, 6, 3, 7, 11, 6, 8, 4, 6, 11, 5, 13, 12, 1},
	{4, 8, 5, 12, 10, 13, 7, 4, 8, 12, 11, 13, 2, 10, 13, 11, 9, 5, 11, 8, 3, 7, 9, 5, 6, 9, 4, 11, 13, 10, 3, 7, 1, 13, 12, 6, 2, 10, 8, 9, 1, 12, 6},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{8, 12, 3, 13, 12, 9, 10, 4, 11, 7, 5, 8, 6, 13, 5, 11, 10, 12, 4, 7, 1, 12, 9, 7, 8, 10, 4, 6, 5, 11, 10, 13, 9, 3, 6, 8, 11, 9, 2, 13},
	{11, 8, 6, 11, 1, 8, 9, 11, 10, 9, 4, 10, 13, 6, 10, 12, 5, 7, 12, 13, 3, 11, 2, 8, 13, 10, 3, 12, 5, 7, 8, 12, 4, 13, 7, 9, 4, 6, 9, 13, 5},
	{12, 13, 6, 3, 13, 2, 9, 4, 11, 10, 13, 9, 5, 6, 10, 4, 8, 5, 7, 10, 3, 11, 12, 7, 11, 5, 9, 12, 1, 13, 7, 11, 12, 8, 13, 6, 10, 8, 4},
	{11, 4, 13, 3, 11, 12, 7, 13, 6, 10, 8, 4, 13, 5, 9, 8, 1, 10, 5, 11, 3, 6, 2, 9, 8, 6, 12, 5, 7, 12, 10, 13, 7, 11, 9, 10, 4, 13, 12},
	{8, 3, 11, 4, 8, 1, 12, 4, 11, 2, 13, 10, 12, 9, 13, 7, 12, 5, 7, 10, 11, 5, 6, 13, 5, 6, 10, 7, 8, 12, 9, 11, 4, 9, 13, 3, 6, 10, 13},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 10, 250, 2500, 10000}, -- wild
	[ 2] = {},                        -- scatter
	[ 3] = {0, 2, 30, 150, 1000},     -- heart
	[ 4] = {0, 2, 20, 125, 500},      -- sword
	[ 5] = {0, 2, 20, 125, 500},      -- shield
	[ 6] = {0, 0, 15, 100, 300},      -- esmerald
	[ 7] = {0, 0, 15, 100, 300},      -- sapphire
	[ 8] = {0, 0, 10, 50, 200},       -- ace
	[ 9] = {0, 0, 10, 50, 200},       -- king
	[10] = {0, 0, 5, 25, 150},        -- queen
	[11] = {0, 0, 5, 25, 150},        -- jack
	[12] = {0, 0, 5, 25, 100},        -- ten
	[13] = {0, 0, 5, 25, 100},        -- nine
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 2, 5, 20, 500}

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 15, 15, 15}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 2 -- minimum scatters to win
local mw = 2 -- multiplier on wilds
local mfs = 3 -- multiplier on free spins

-- Performs full RTP calculation for given reels
local function calculate(reels_reg, reels_bon)
	assert(#reels_reg == sx, "unexpected number of regular reels")
	assert(#reels_bon == sx, "unexpected number of bonus reels")

	local reels
	local N, L, counts

	-- Reels precalculations
	local function precalculate_reels()
		-- Get number of total reshuffles and lengths of each reel.
		N, L = 1, {}
		for i, r in ipairs(reels) do
			N = N * #r
			L[i] = #r
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
						-- Total combinations
						local combs_total = 1
						for i = 1, sx do
							if i <= n then
								combs_total = combs_total * c[i]
							elseif i == n + 1 then
								combs_total = combs_total * (L[i] - c[i])
							else
								combs_total = combs_total * L[i]
							end
						end

						-- Combinations WITHOUT any wilds on reels
						local combs_no_wild = 1
						for i = 1, sx do
							if i <= n then
								combs_no_wild = combs_no_wild * s[i]
							elseif i == n + 1 then
								combs_no_wild = combs_no_wild * (L[i] - c[i])
							else
								combs_no_wild = combs_no_wild * L[i]
							end
						end

						local better_wilds = 0
						local wn_min = nil
						for wn = line_min, n do
							if wpays[wn] >= payout*mw then
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
									bw = bw * (L[i] - c[i])
								else
									bw = bw * L[i]
								end
							end
							better_wilds = bw
						end

						local combs_with_wild = combs_total - combs_no_wild - better_wilds
						ev_sum = ev_sum + (combs_no_wild + combs_with_wild * mw) * payout
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
					elseif i == n + 1 then wc = wc * (L[i] - w[i])
					else wc = wc * L[i] end
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
								if pays[sn]*mw > payout then
									local loss = 1
									for i = 1, sx do
										if i <= n then
											loss = loss * w[i]
										elseif i == n + 1 then
											loss = loss * s[i]
										elseif i <= sn then
											loss = loss * c[i]
										elseif i == sn + 1 then
											loss = loss * (L[i] - c[i])
										else
											loss = loss * L[i]
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
				current_comb * (L[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

		return ev_sum, fs_sum, fs_num
	end

	-- Execute calculation
	local rtp_fs
	do
		reels = reels_bon
		precalculate_reels()
		local rtp_line = calculate_line_ev() / N
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / N
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / N
		local sq = 1 / (1 - q)
		rtp_fs = mfs * sq * rtp_sym
		print(string.format("*bonus reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", N/fs_num))
		print(string.format("RTP = %g*sq*rtp(sym) = %g*%.5g*%.5g = %.6f%%", mfs, mfs, sq, rtp_sym*100, rtp_fs*100))
	end
	local rtp_total
	do
		reels = reels_reg
		precalculate_reels()
		local rtp_line = calculate_line_ev() / N
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / N
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / N
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", N/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym*100, q, rtp_fs*100, rtp_total*100))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
