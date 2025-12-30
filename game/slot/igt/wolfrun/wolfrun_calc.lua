-- IGT / Wolf Run
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{1, 1, 1, 1, 3, 8, 5, 7, 9, 4, 11, 2, 8, 4, 10, 5, 8, 11, 9, 5, 6, 2, 7, 10, 11, 3, 6, 7, 4, 10, 3, 9, 2, 6},
	{1, 1, 1, 1, 10, 4, 8, 11, 6, 4, 7, 5, 9, 8, 3, 11, 12, 5, 7, 2, 6, 10, 9, 3, 12, 4, 10, 2, 11, 9, 3, 7, 6, 8, 2, 12, 5},
	{1, 1, 1, 1, 8, 5, 10, 3, 11, 6, 9, 12, 2, 7, 6, 3, 11, 4, 7, 2, 12, 4, 8, 10, 5, 6, 4, 9, 3, 10, 8, 5, 7, 12, 9, 2, 11},
	{1, 1, 1, 1, 1, 3, 9, 12, 10, 4, 8, 11, 5, 6, 4, 7, 9, 3, 10, 4, 11, 8, 5, 7, 2, 9, 12, 7, 3, 6, 12, 2, 8, 5, 11, 2, 6, 10},
	{1, 1, 1, 1, 1, 1, 11, 2, 6, 9, 5, 11, 8, 4, 10, 5, 6, 8, 2, 7, 9, 3, 6, 5, 11, 3, 7, 4, 10, 2, 7, 9, 3, 8, 4, 10},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{1, 1, 1, 1, 2, 11, 10, 8, 4, 11, 10, 9, 7, 6, 8, 9, 11, 6, 3, 7, 5, 8, 10, 6, 2, 9, 7, 5, 10, 11, 2, 8, 7, 4, 9, 3, 7, 10, 6, 5, 11, 9, 10, 8, 3, 11, 6, 8, 4, 9},
	{1, 1, 1, 1, 12, 11, 10, 9, 3, 8, 2, 6, 9, 7, 10, 8, 6, 7, 2, 10, 5, 9, 8, 11, 4, 7, 9, 11, 12, 5, 10, 6, 8, 12, 9, 7, 8, 11, 10, 4, 7, 3, 11, 4, 8, 6, 5, 11, 10, 6, 3, 9, 2},
	{1, 1, 1, 1, 4, 10, 6, 11, 8, 9, 4, 7, 11, 9, 3, 7, 8, 11, 2, 9, 7, 10, 6, 5, 8, 11, 2, 10, 6, 9, 7, 12, 10, 5, 11, 7, 4, 9, 8, 10, 3, 12, 2, 9, 8, 6, 11, 5, 8, 10, 6, 3, 12},
	{1, 1, 1, 1, 1, 9, 3, 12, 11, 9, 6, 4, 10, 2, 8, 9, 11, 6, 8, 5, 7, 9, 11, 2, 10, 3, 6, 12, 5, 10, 2, 8, 7, 6, 10, 9, 4, 8, 6, 7, 11, 5, 9, 8, 3, 10, 12, 11, 7, 4, 8, 11, 10, 7},
	{1, 1, 1, 1, 11, 9, 4, 10, 3, 8, 9, 10, 6, 11, 5, 10, 8, 9, 2, 7, 8, 11, 6, 9, 7, 5, 6, 3, 9, 10, 7, 11, 5, 10, 7, 4, 6, 8, 11, 3, 10, 2, 9, 4, 6, 8, 11, 7, 2, 8},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 50, 200, 1000}, -- moon wolf
	[ 2] = {0, 0, 25, 100, 400},  -- grey wolf
	[ 3] = {0, 0, 25, 100, 400},  -- white wolf
	[ 4] = {0, 0, 20, 75, 250},   -- idol1
	[ 5] = {0, 0, 20, 75, 250},   -- idol2
	[ 6] = {0, 0, 5, 50, 150},    -- ace
	[ 7] = {0, 0, 5, 50, 150},    -- king
	[ 8] = {0, 0, 5, 20, 100},    -- queen
	[ 9] = {0, 0, 5, 20, 100},    -- jack
	[10] = {0, 0, 5, 20, 100},    -- ten
	[11] = {0, 0, 5, 20, 100},    -- nine
	[12] = {},                    -- bonus (2, 3, 4 reels only)
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_pay, scat_fs = 2, 5 -- scatter pays and number of free spins awarded

-- 4. CONFIGURATION
local sx, sy = 5, 4 -- screen width & height
local wild, scat = 1, 12 -- wild & scatter symbol IDs
local line_min = 3 -- minimum line symbols to win
local scat_min = 3 -- minimum scatters to win

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
		for symbol_id in pairs(PAYTABLE_LINE) do
			counts[symbol_id] = {}
			for i = 1, sx do counts[symbol_id][i] = 0 end
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
		rtp_fs = sq * rtp_sym
		print(string.format("*bonus reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%", sq, rtp_sym, rtp_fs))
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
		print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
