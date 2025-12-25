-- AGT / Wizard
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{10, 10, 10, 10, 6, 8, 8, 8, 8, 9, 9, 9, 9, 4, 4, 4, 4, 5, 5, 5, 5, 3, 7, 7, 7, 7, 1, 1, 1, 1, 7, 8, 10, 10, 11, 11, 1, 1, 1, 3, 3, 3, 3, 5, 6, 6, 6, 6, 9, 9, 4, 11, 11, 11, 11},
	{11, 10, 10, 10, 10, 6, 11, 11, 11, 11, 9, 4, 4, 4, 4, 8, 3, 3, 3, 3, 1, 1, 1, 5, 9, 9, 9, 9, 4, 7, 7, 7, 7, 3, 6, 6, 6, 6, 8, 8, 8, 8, 5, 5, 5, 5, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 10, 7},
	{10, 10, 10, 10, 6, 4, 6, 6, 6, 6, 1, 1, 1, 8, 5, 5, 5, 5, 9, 8, 8, 8, 8, 11, 4, 4, 4, 4, 3, 7, 7, 7, 7, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 10, 7, 1, 1, 1, 1, 5, 9, 9, 9, 9, 3, 3, 3, 3, 11, 11, 11, 11},
	{3, 5, 5, 5, 5, 8, 8, 8, 8, 5, 6, 6, 6, 6, 7, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 7, 7, 7, 7, 10, 4, 10, 10, 10, 10, 11, 11, 11, 11, 9, 9, 9, 9, 4, 4, 4, 4, 9, 11, 8, 6, 1, 1, 1, 1, 3, 3, 3, 3, 1, 1, 1},
	{10, 10, 3, 7, 10, 10, 10, 10, 5, 6, 6, 6, 6, 7, 7, 7, 7, 1, 1, 1, 8, 9, 9, 3, 3, 3, 3, 4, 8, 8, 8, 8, 5, 5, 5, 5, 11, 11, 6, 4, 4, 4, 4, 9, 9, 9, 9, 1, 1, 1, 1, 11, 11, 11, 11},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 10, 50, 250, 1000}, -- wild
	[ 2] = {0, 0, 0, 0, 0},        -- scatter (2, 3, 4 reels only)
	[ 3] = {0, 5, 10, 20, 50},     -- owl
	[ 4] = {0, 4, 10, 20, 45},     -- cat
	[ 5] = {0, 0, 8, 15, 40},      -- cauldron
	[ 6] = {0, 0, 8, 15, 30},      -- emerald
	[ 7] = {0, 0, 6, 15, 25},      -- ruby
	[ 8] = {0, 0, 4, 10, 20},      -- ace
	[ 9] = {0, 0, 4, 10, 20},      -- king
	[10] = {0, 0, 2, 8, 15},       -- queen
	[11] = {0, 0, 2, 8, 15},       -- jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay = 2

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 15, 30}

-- 5. CONFIGURATION
local sx, sy = 5, 4 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local line_min = 2 -- minimum line symbols to win
local scat_min = 10 -- minimum scatters to win

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
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- 1. Preliminary calculation: Determine how many ways
		-- N scatter symbols can appear on each reel.
		-- ways_on_reel[reel_idx][scat_count] = number of combinations
		-- that yield scat_count scatters
		local ways_on_reel = {}
		for i, reel in ipairs(reels) do
			ways_on_reel[i] = {}
			for count = 0, sy do
				ways_on_reel[i][count] = 0
			end
			for stop_idx = 0, lens[i] - 1 do
				local count = 0
				for h = 0, sy - 1 do
					local symbol_idx = (stop_idx + h) % lens[i] + 1
					if reel[symbol_idx] == scat then
						count = count + 1
					end
				end
				ways_on_reel[i][count] = ways_on_reel[i][count] + 1
			end
		end

		-- 2. Recursive traversal of the combination tree to sum the EV
		-- reel_idx: current reel (1-5)
		-- scat_sum: the sum of scatters on the previous reels
		-- current_comb: the number of combinations that lead to this state
		local function find_scatter_ev_recursive(reel_idx, scat_sum, current_comb)

			-- Base case: all 5 reels have been processed
			if reel_idx > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * scat_pay
					fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
					fs_num = fs_num + current_comb
				end
				return
			end

			-- Recursive step: iterate through all possible outcomes
			-- of the current reel (0, 1, 2, or 3 scatters)
			for scat_count, ways in pairs(ways_on_reel[reel_idx]) do
				if ways > 0 then
					find_scatter_ev_recursive(
						reel_idx + 1,
						scat_sum + scat_count,
						current_comb * ways
					)
				end
			end
		end
		find_scatter_ev_recursive(1, 0, 1)

		return ev_sum, fs_sum, fs_num
	end

	-- Execute calculation
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles * 100
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = sq * rtp_sym
	local rtp_total = rtp_sym + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
