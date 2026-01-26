-- CT Interactive / Alaska Wild
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{1, 1, 1, 1, 5, 7, 12, 6, 9, 11, 8, 6, 12, 4, 10, 6, 12, 5, 10, 4, 3, 9, 12, 4, 3, 11, 5, 8, 7, 11, 12, 9, 7, 11, 8, 9, 3, 10, 6, 8, 5, 10, 7, 11, 4, 10, 9, 3},
	{1, 1, 1, 1, 8, 9, 10, 7, 4, 3, 10, 11, 9, 5, 2, 12, 4, 8, 3, 6, 11, 12, 9, 8, 11, 5, 10, 12, 4, 9, 7, 8, 6, 2, 9, 3, 6, 4, 7, 2, 10, 12, 11, 7, 5, 12, 8, 3, 6, 10, 11, 5},
	{1, 1, 1, 1, 5, 2, 11, 8, 9, 7, 10, 4, 11, 3, 9, 8, 2, 6, 12, 4, 3, 7, 10, 12, 11, 7, 8, 4, 11, 10, 5, 9, 6, 11, 12, 7, 5, 10, 8, 9, 6, 5, 3, 8, 6, 12, 10, 2, 4, 9, 3, 12},
	{1, 1, 1, 1, 7, 12, 6, 8, 4, 7, 10, 9, 6, 8, 11, 9, 5, 3, 8, 12, 4, 7, 5, 6, 9, 11, 10, 12, 3, 6, 9, 10, 8, 2, 3, 5, 11, 4, 12, 3, 8, 11, 10, 2, 12, 4, 7, 10, 2, 5, 11, 9},
	{1, 1, 1, 1, 3, 12, 5, 6, 8, 9, 7, 3, 12, 4, 8, 6, 10, 5, 7, 9, 4, 3, 11, 10, 6, 3, 5, 9, 11, 12, 10, 8, 9, 7, 4, 11, 10, 8, 12, 6, 9, 11, 7, 12, 4, 5, 10, 11},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 50, 200, 1000}, --  1 wild
	[ 2] = {},                    --  2 scatter (2, 3, 4 reels only)
	[ 3] = {0, 0, 25, 100, 400},  --  3 fox
	[ 4] = {0, 0, 25, 100, 400},  --  4 squirrel
	[ 5] = {0, 0, 15, 75, 250},   --  5 owl
	[ 6] = {0, 0, 15, 75, 250},   --  6 eagle
	[ 7] = {0, 0, 5, 50, 150},    --  7 rockfish
	[ 8] = {0, 0, 5, 50, 150},    --  8 bulltrout
	[ 9] = {0, 0, 5, 15, 100},    --  9 ace
	[10] = {0, 0, 5, 15, 100},    -- 10 king
	[11] = {0, 0, 5, 15, 100},    -- 11 queen
	[12] = {0, 0, 5, 15, 100},    -- 12 jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay, scat_fs = 2, 5 -- scatter pays and number of free spins awarded

-- 4. CONFIGURATION
local sx, sy = 5, 4 -- screen width & height
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
					ev_sum = ev_sum + current_comb * scat_pay
					fs_sum = fs_sum + current_comb * scat_fs
					fs_num = fs_num + current_comb
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
