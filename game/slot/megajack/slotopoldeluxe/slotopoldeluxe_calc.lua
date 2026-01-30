-- Megajack / Slotopol Deluxe
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{1, 5, 13, 4, 13, 3, 13, 5, 9, 7, 8, 13, 10, 13, 12, 11, 13, 12, 11, 13, 13, 2, 4, 5, 2, 6, 7, 9, 8, 3, 10, 6},
	{13, 2, 12, 9, 13, 4, 5, 6, 9, 7, 13, 10, 12, 13, 11, 13, 13, 11, 12, 13, 3, 4, 5, 2, 8, 7, 10, 4, 6, 8, 3, 1},
	{2, 1, 12, 3, 4, 5, 2, 6, 7, 10, 8, 4, 5, 13, 12, 11, 13, 12, 11, 13, 12, 3, 5, 13, 9, 6, 7, 10, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 6, 4, 2, 7, 8, 3, 5, 13, 12, 11, 13, 12, 11, 13, 12, 2, 4, 5, 3, 12, 6, 10, 7, 13, 8, 9, 13},
	{1, 2, 12, 4, 3, 5, 12, 6, 7, 3, 8, 12, 2, 13, 12, 11, 13, 12, 11, 13, 12, 3, 4, 5, 2, 6, 7, 10, 13, 8, 9, 5},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                        -- dollar
	[ 2] = {0, 2, 5, 25, 100},        -- cherry
	[ 3] = {0, 2, 5, 25, 100},        -- plum
	[ 4] = {0, 0, 5, 25, 100},        -- wmelon
	[ 5] = {0, 0, 5, 25, 100},        -- grapes
	[ 6] = {0, 0, 10, 100, 250},      -- ananas
	[ 7] = {0, 0, 10, 100, 250},      -- lemon
	[ 8] = {0, 0, 10, 100, 250},      -- drink
	[ 9] = {0, 2, 10, 100, 500},      -- palm
	[10] = {0, 2, 10, 100, 500},      -- yacht
	[11] = {0, 10, 200, 2000, 10000}, -- eldorado
	[12] = {},                        -- spin
	[13] = {},                        -- dice
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 2, 20, 1000}

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 11, 1 -- wild & scatter symbol IDs
local bon1, bon2 = 12, 13 -- bonus games symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 3 -- minimum scatters to win
local mw = 2 -- multiplier on wilds
local EVmje1 = 106 * 1 -- Eldorado 1 spins bonus expectation
local EVmje3 = 106 * 3 -- Eldorado 3 spins bonus expectation
local EVmje6 = 106 * 6 -- Eldorado 6 spins bonus expectation
local EVmjm = 286.6059742226795 -- Monopoly bonus expectation

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
									bw = bw * (lens[i] - c[i])
								else
									bw = bw * lens[i]
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
				current_comb * (lens[reel_index] - c[reel_index] * sy))
		end
		find_scatter_combs(1, 0, 1) -- Start recursion

		return ev_sum
	end

	-- Calculating Eldorado1 bonus symbols
	local function calculate_mje1_comb()
		local b = counts[bon1]
		local comb5 = b[1] * b[2] * b[3] * (lens[4] - b[4]) * lens[5]
		return comb5
	end

	-- Calculating Eldorado3 bonus symbols
	local function calculate_mje3_comb()
		local b = counts[bon1]
		local comb5 = b[1] * b[2] * b[3] * b[4] * (lens[5] - b[5])
		return comb5
	end

	-- Calculating Eldorado6 bonus symbols
	local function calculate_mje6_comb()
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
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local rtp_scat = calculate_scat_ev() / reshuffles * 100
	local rtp_sym = rtp_line + rtp_scat
	local qmje1 = calculate_mje1_comb() / reshuffles
	local rtp_mje1 = EVmje1 * qmje1 * 100
	local qmje3 = calculate_mje3_comb() / reshuffles
	local rtp_mje3 = EVmje3 * qmje3 * 100
	local qmje6 = calculate_mje6_comb() / reshuffles
	local rtp_mje6 = EVmje6 * qmje6 * 100
	local comb_mjm = calculate_mjm_comb()
	local qmjm = comb_mjm / reshuffles
	local rtp_mjm = EVmjm * qmjm * 100
	local rtp_total = rtp_sym + rtp_mje1 + rtp_mje3 + rtp_mje6 + rtp_mjm
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
	print(string.format("spin1 bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmje1, rtp_mje1))
	print(string.format("spin3 bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmje3, rtp_mje3))
	print(string.format("spin6 bonuses: hit rate 1/%.5g, rtp = %.6f%%", 1/qmje6, rtp_mje6))
	print(string.format("monopoly bonuses: hit rate 1/%.5g, rtp = %.6f%%", reshuffles/comb_mjm, rtp_mjm))
	print(string.format("RTP = %.5g(sym) + %.5g(mje) + %.5g(mjm) = %.6f%%",
		rtp_sym, rtp_mje1 + rtp_mje3 + rtp_mje6, rtp_mjm, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
