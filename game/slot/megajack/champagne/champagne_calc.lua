-- Megajack / Champagne
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{12, 1, 5, 2, 12, 11, 2, 11, 12, 3, 2, 8, 12, 3, 4, 6, 12, 2, 5, 10, 3, 9, 7, 8, 4, 3, 7, 9, 2, 3, 4, 6},
	{2, 5, 10, 12, 9, 6, 3, 4, 12, 2, 6, 8, 3, 12, 11, 2, 11, 12, 5, 7, 4, 6, 3, 4, 12, 2, 5, 8, 2, 7, 1, 9},
	{12, 5, 10, 12, 9, 6, 3, 4, 12, 2, 6, 8, 3, 12, 11, 2, 11, 12, 5, 7, 4, 6, 3, 4, 12, 2, 5, 8, 12, 7, 1, 9},
	{12, 8, 2, 12, 6, 5, 2, 4, 12, 2, 1, 3, 2, 9, 7, 12, 11, 11, 11, 11, 12, 5, 2, 12, 8, 6, 2, 3, 10, 12, 2, 4},
	{12, 11, 7, 12, 6, 4, 12, 3, 2, 12, 3, 7, 12, 3, 5, 1, 12, 3, 8, 9, 12, 4, 3, 2, 12, 5, 3, 10, 2, 12, 3, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                        -- dollar
	[ 2] = {0, 3, 5, 20, 100},        -- cherry
	[ 3] = {0, 3, 5, 20, 100},        -- plum
	[ 4] = {0, 0, 5, 20, 100},        -- wmelon
	[ 5] = {0, 0, 5, 20, 100},        -- grapes
	[ 6] = {0, 0, 5, 20, 100},        -- ananas
	[ 7] = {0, 0, 5, 20, 100},        -- lemon
	[ 8] = {0, 0, 5, 20, 100},        -- drink
	[ 9] = {0, 5, 10, 20, 1000},      -- palm
	[10] = {0, 7, 10, 20, 1000},      -- yacht
	[11] = {0, 10, 100, 2000, 10000}, -- eldorado
	[12] = {},                        -- fizz
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 0, 0, 1000}

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 15, 15, 15}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat, bon = 11, 1, 12 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 2 -- minimum scatters to win
local mfs = 2 -- multiplier on free spins
local EVmjc = 168 -- Bottle game calculated expectation

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
	local function calculate_line_ev(free_spins)
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

						local mm = free_spins and mfs or 1
						ev_sum = ev_sum + (combs_total - better_wilds) * payout * mm
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
				local mm = (free_spins and n < sx) and mfs or 1
				ev_sum = ev_sum + (wc - losses) * payout * mm
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
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
					if FREESPIN_SCAT[scat_sum] > 0 then
						fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum] * mm
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

	-- Calculating Monopoly bonus symbols
	local function calculate_mjc_comb()
		local b = counts[bon]
		local comb5 = b[1] * b[2] * b[3] * b[4] * b[5]
		return comb5
	end

	-- Execute calculation
	local rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	do
		local rtp_line = calculate_line_ev(true) / reshuffles * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(true)
		local rtp_scat = ev_sum / reshuffles * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		local qmjc = calculate_mjc_comb() / reshuffles
		local rtp_mjc = EVmjc * qmjc * 100
		rtp_fs = sq * (rtp_sym + rtp_mjc)
		print(string.format("*free games calculations*"))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("bottle bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmjc, rtp_mjc))
		print(string.format("RTP = sq*(rtp(sym)+rtp(mjc)) = %.5g*(%.5g+%.5g) = %.6f%%", sq, rtp_sym, rtp_mjc, rtp_fs))
	end
	local rtp_total
	do
		local rtp_line = calculate_line_ev(false) / reshuffles * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(false)
		local rtp_scat = ev_sum / reshuffles * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		local qmjc = calculate_mjc_comb() / reshuffles
		local rtp_mjc = EVmjc * qmjc * 100
		rtp_total = rtp_sym + rtp_mjc + q * rtp_fs
		print(string.format("*regular games calculations*"))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("bottle bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmjc, rtp_mjc))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
