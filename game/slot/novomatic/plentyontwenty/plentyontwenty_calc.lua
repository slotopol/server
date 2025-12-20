-- Novomatic / Plenty on Twenty
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{8, 8, 8, 3, 4, 4, 4, 4, 5, 5, 5, 5, 6, 6, 6, 5, 5, 5, 5, 8, 8, 8, 8, 8, 8, 3, 3, 3, 3, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 2, 6, 6, 6, 6, 6, 6, 8, 8, 8, 8, 8, 8, 1, 1, 1, 4, 4, 4, 4, 6, 6, 6, 6, 6, 6, 4, 4, 4, 4, 2, 7, 7, 7, 7, 7, 7, 2, 7, 7, 7, 3, 3, 3, 3},
	{3, 4, 4, 4, 4, 5, 5, 5, 5, 4, 4, 4, 4, 6, 6, 6, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 5, 5, 5, 5, 8, 8, 8, 3, 3, 3, 3, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 8, 8, 8, 8, 8, 8, 2, 7, 7, 7, 4, 4, 4, 4, 5, 5, 5, 5, 1, 1, 1, 3, 3, 3, 3, 8, 8, 8, 8, 8, 8, 2},
	{5, 5, 5, 5, 2, 6, 6, 6, 6, 6, 6, 7, 7, 7, 8, 8, 8, 8, 8, 8, 3, 7, 7, 7, 7, 7, 7, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 4, 4, 4, 4, 8, 8, 8, 8, 8, 8, 2, 3, 3, 3, 3, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6, 3, 3, 3, 3, 6, 6, 6, 4, 4, 4, 4, 5, 5, 5, 5, 2, 8, 8, 8, 1, 1, 1},
	{6, 6, 6, 3, 3, 3, 3, 6, 6, 6, 6, 6, 6, 5, 5, 5, 5, 6, 6, 6, 6, 6, 6, 3, 3, 3, 3, 2, 3, 1, 1, 1, 4, 4, 4, 4, 8, 8, 8, 4, 4, 4, 4, 8, 8, 8, 8, 8, 8, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 8, 8, 8, 8, 8, 8, 2, 4, 4, 4, 4, 7, 7, 7, 7, 7, 7, 5, 5, 5, 5, 7, 7, 7},
	{7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 8, 8, 8, 8, 8, 8, 6, 6, 6, 6, 6, 6, 4, 4, 4, 4, 8, 8, 8, 8, 8, 8, 5, 5, 5, 5, 3, 5, 5, 5, 5, 2, 1, 1, 1, 3, 3, 3, 3, 7, 7, 7, 7, 7, 7, 3, 3, 3, 3, 5, 5, 5, 5, 2, 4, 4, 4, 4, 7, 7, 7, 8, 8, 8, 2, 4, 4, 4, 4, 6, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[1] = {0, 0, 40, 400, 1000}, -- wild
	[2] = {0, 0, 0, 0, 0},       -- scatter
	[3] = {0, 0, 20, 80, 400},   -- bell
	[4] = {0, 0, 20, 40, 200},   -- melon
	[5] = {0, 0, 20, 40, 200},   -- plum
	[6] = {0, 0, 10, 20, 100},   -- orange
	[7] = {0, 0, 10, 20, 100},   -- lemon
	[8] = {0, 0, 10, 20, 100},   -- cherry
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 5, 20, 500}

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local line_min = 3 -- minimum line symbols to win
local scat_min = 3 -- minimum scatters to win

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
		local wpays = PAYTABLE_LINE[wild]

		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			if symbol_id ~= scat and symbol_id ~= wild then
				local s = counts[symbol_id]
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end

				for n = line_min, sx do
					local payout = pays[n]
					if payout > 0 then
						local total_combs = 1
						for i = 1, sx do
							if i <= n then
								total_combs = total_combs * c[i]
							elseif i == n + 1 then
								total_combs = total_combs * (lens[i] - c[i])
							else
								total_combs = total_combs * lens[i]
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

						ev_sum = ev_sum + (total_combs - better_wilds) * payout
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
					for symbol_id, pays in pairs(PAYTABLE_LINE) do
						if symbol_id ~= wild and symbol_id ~= scat then
							local s = counts[symbol_id]
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

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local rtp_scat = calculate_scat_ev() / reshuffles * 100
	local rtp_total = rtp_line + rtp_scat
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("RTP = %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
