-- Megajack / Slotopol
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{13, 1, 5, 12, 13, 11, 12, 11, 13, 8, 2, 12, 13, 3, 4, 6, 13, 2, 5, 10, 13, 9, 7, 8, 13, 10, 7, 9, 13, 3, 4, 6},
	{9, 5, 10, 13, 9, 6, 3, 4, 13, 2, 12, 8, 12, 13, 11, 12, 11, 13, 5, 7, 10, 6, 3, 4, 13, 2, 12, 8, 13, 7, 1, 12},
	{12, 13, 11, 12, 11, 13, 5, 10, 9, 7, 1, 12, 13, 3, 8, 6, 12, 13, 8, 4, 12, 2, 5, 10, 13, 7, 2, 13, 6, 3, 4, 9},
	{12, 1, 2, 13, 6, 5, 12, 4, 8, 12, 13, 3, 10, 9, 7, 13, 11, 11, 11, 11, 13, 5, 12, 9, 8, 6, 13, 3, 10, 2, 7, 4},
	{13, 11, 13, 12, 6, 4, 12, 3, 2, 5, 12, 10, 7, 12, 8, 1, 9, 12, 8, 9, 12, 4, 3, 12, 2, 5, 12, 10, 7, 13, 12, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                        -- dollar
	[ 2] = {0, 2, 5, 15, 100},        -- cherry
	[ 3] = {0, 2, 5, 15, 100},        -- plum
	[ 4] = {0, 0, 5, 15, 100},        -- wmelon
	[ 5] = {0, 0, 5, 15, 100},        -- grapes
	[ 6] = {0, 0, 5, 15, 100},        -- ananas
	[ 7] = {0, 0, 5, 15, 100},        -- lemon
	[ 8] = {0, 0, 5, 15, 100},        -- drink
	[ 9] = {0, 2, 5, 15, 100},        -- palm
	[10] = {0, 2, 5, 15, 100},        -- yacht
	[11] = {0, 10, 100, 2000, 10000}, -- eldorado
	[12] = {},                        -- spin
	[13] = {},                        -- dice
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 5, 8, 20, 1000}

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 11, 1 -- wild & scatter symbol IDs
local bon1, bon2 = 12, 13 -- bonus games symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 2 -- minimum scatters to win
local mw = 2 -- multiplier on wilds
local EVmje9 = 106 * 9 -- Eldorado 9 spins bonus expectation
local EVmjm = 286.6059742226795 -- Monopoly bonus expectation

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
		local ev_sum = 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
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

		return ev_sum
	end

	-- Calculating Eldorado9 bonus symbols
	local function calculate_mje9_comb()
		local b = counts[bon1]
		local comb5 = b[1] * b[2] * b[3] * b[4] * b[5]
		return comb5
	end

	-- Calculating Monopoly bonus symbols
	local function calculate_mjm_comb()
		local b = counts[bon2]
		local comb5 = b[1] * b[2] * b[3] * b[4] * b[5]
		return comb5
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / N
	local rtp_scat = calculate_scat_ev() / N
	local rtp_sym = rtp_line + rtp_scat
	local qmje9 = calculate_mje9_comb() / N
	local rtp_mje9 = EVmje9 * qmje9
	local qmjm = calculate_mjm_comb() / N
	local rtp_mjm = EVmjm * qmjm
	local rtp_total = rtp_sym + rtp_mje9 + rtp_mjm
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(L, ", "), N))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line*100, rtp_scat*100, rtp_sym*100))
	print(string.format("spin9 bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmje9, rtp_mje9*100))
	print(string.format("monopoly bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmjm, rtp_mjm*100))
	print(string.format("RTP = %.5g(sym) + %.5g(mje9) + %.5g(mjm) = %.6f%%",
		rtp_sym*100, rtp_mje9*100, rtp_mjm*100, rtp_total*100))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
