-- AGT / Egypt
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{8, 5, 4, 9, 11, 2, 7, 11, 1, 6, 9, 5, 10, 11, 9, 8, 7, 12, 1, 6, 12, 5, 3, 10, 8, 9, 11, 8, 12, 10, 11, 12, 10, 8, 11, 7, 12, 3, 6, 4, 12, 7, 10, 9},
	{4, 11, 12, 5, 11, 9, 10, 6, 12, 9, 2, 11, 10, 7, 1, 11, 12, 9, 6, 8, 5, 11, 7, 8, 10, 7, 11, 9, 4, 10, 3, 7, 10, 12, 9, 8, 6, 3, 8, 12, 5, 8, 12, 1},
	{12, 9, 11, 1, 7, 10, 11, 9, 7, 10, 11, 4, 12, 5, 10, 11, 6, 12, 5, 6, 8, 2, 5, 8, 7, 11, 8, 10, 1, 9, 12, 10, 9, 8, 11, 6, 12, 3, 8, 7, 12, 9, 4, 3},
	{11, 9, 12, 10, 8, 6, 4, 9, 5, 6, 2, 12, 3, 5, 8, 12, 11, 7, 9, 6, 11, 12, 1, 8, 11, 10, 4, 11, 7, 9, 12, 3, 10, 8, 11, 7, 5, 1, 12, 10, 7, 8, 9, 10},
	{11, 9, 8, 6, 11, 2, 9, 6, 8, 7, 12, 10, 3, 9, 12, 8, 11, 12, 10, 9, 7, 11, 6, 12, 11, 10, 5, 3, 12, 7, 9, 10, 11, 4, 1, 12, 8, 1, 5, 4, 8, 10, 5, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 10, 250, 2500, 9000}, -- wild
	[ 2] = {},                       -- scatter
	[ 3] = {0, 2, 25, 125, 750},     -- anubis
	[ 4] = {0, 2, 25, 125, 750},     -- sphinx
	[ 5] = {0, 0, 20, 100, 400},     -- snake
	[ 6] = {0, 0, 15, 75, 250},      -- mummy
	[ 7] = {0, 0, 15, 75, 250},      -- wall
	[ 8] = {0, 0, 10, 50, 125},      -- cat
	[ 9] = {0, 0, 10, 50, 125},      -- ace
	[10] = {0, 0, 5, 25, 100},       -- king
	[11] = {0, 0, 5, 10, 50},        -- queen
	[12] = {0, 0, 5, 10, 25},        -- jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 2, 5, 20, 500}

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 2 -- minimum scatters to win
local EVmini = (24*1 + 1*3 + 1*6 + 1*9)/27 -- expectation value on minislot

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
			if symbol_id ~= wild and #pays > 0 then
				local s = counts[symbol_id]
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
					for symbol_id, pays in pairs(PAYTABLE_LINE) do
						if symbol_id ~= wild and #pays > 0 then
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
	local rtp_sym = rtp_line + rtp_scat
	local rtp_total = rtp_sym * EVmini
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
	print(string.format("RTP = %.5g(sym) * %.5g(Em) = %.6f%%", rtp_sym, EVmini, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
