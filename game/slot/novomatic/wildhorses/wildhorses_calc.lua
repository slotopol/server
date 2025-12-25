-- Novomatic / Wild Horses
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{1, 1, 1, 1, 7, 7, 11, 8, 9, 3, 8, 5, 7, 4, 11, 7, 10, 8, 9, 4, 7, 3, 10, 12, 8, 4, 9, 10, 6, 8, 9, 12, 10, 8, 5, 12, 12, 10, 6, 11, 12, 9, 6, 11, 3, 9, 11, 11, 10, 5, 7, 12},
	{1, 1, 1, 1, 8, 7, 9, 4, 10, 9, 11, 12, 8, 3, 9, 9, 6, 7, 10, 8, 9, 11, 7, 7, 4, 10, 8, 11, 7, 6, 8, 12, 11, 12, 4, 12, 10, 5, 12, 3, 11, 5, 11, 6, 12, 10, 7, 5, 8, 10, 3, 9},
	{1, 1, 1, 1, 8, 7, 8, 3, 12, 7, 12, 12, 5, 11, 9, 8, 3, 10, 10, 6, 9, 5, 11, 10, 12, 8, 11, 8, 4, 11, 8, 3, 7, 10, 11, 10, 6, 9, 9, 7, 9, 5, 12, 7, 10, 4, 9, 7, 12, 4, 11, 6},
	{1, 1, 1, 1, 5, 7, 10, 9, 6, 12, 9, 3, 9, 11, 10, 8, 8, 6, 8, 12, 5, 12, 8, 9, 3, 12, 4, 11, 10, 8, 11, 10, 4, 11, 7, 6, 8, 10, 9, 9, 7, 11, 3, 12, 11, 5, 7, 10, 4, 7, 12, 7},
	{1, 1, 1, 1, 8, 12, 11, 10, 5, 10, 4, 7, 11, 11, 9, 10, 9, 12, 12, 8, 7, 7, 11, 4, 7, 3, 8, 4, 9, 8, 9, 6, 7, 11, 8, 10, 6, 9, 11, 10, 6, 10, 5, 12, 7, 3, 8, 9, 5, 12, 3, 12},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{1, 1, 1, 1, 7, 5, 10, 12, 8, 4, 10, 5, 7, 11, 4, 9, 12, 3, 9, 12, 5, 11, 10, 3, 12, 12, 12, 4, 12, 9, 6, 12, 4, 7, 11, 6, 10, 12, 8, 8, 9, 8, 10, 7, 8, 6, 8, 11, 8, 6, 11, 11, 9, 7, 3, 9, 5, 7, 10, 10, 3, 11, 11, 7, 9, 8, 10, 9},
	{1, 1, 1, 1, 4, 8, 10, 11, 5, 7, 10, 12, 5, 10, 6, 10, 9, 9, 12, 12, 3, 8, 5, 10, 7, 9, 7, 3, 12, 12, 4, 8, 11, 6, 8, 9, 2, 7, 3, 11, 7, 12, 11, 12, 6, 7, 2, 11, 8, 10, 2, 8, 10, 9, 10, 4, 9, 3, 9, 7, 4, 11, 11, 11, 12, 8, 12, 6, 8, 5, 9},
	{1, 1, 1, 1, 9, 10, 9, 8, 12, 5, 10, 6, 7, 11, 4, 12, 6, 7, 4, 9, 7, 10, 10, 11, 10, 3, 8, 11, 2, 8, 5, 11, 6, 12, 3, 8, 12, 5, 9, 11, 8, 12, 11, 5, 12, 7, 4, 8, 2, 7, 8, 10, 9, 4, 7, 8, 10, 3, 12, 10, 9, 3, 12, 6, 9, 2, 9, 11, 7, 7, 11},
	{1, 1, 1, 1, 12, 11, 11, 12, 5, 12, 3, 9, 6, 11, 2, 9, 8, 6, 9, 12, 4, 7, 3, 7, 10, 4, 8, 7, 10, 7, 9, 10, 7, 8, 3, 12, 4, 8, 9, 9, 6, 11, 11, 11, 8, 10, 7, 7, 6, 8, 5, 8, 2, 7, 12, 4, 8, 5, 10, 9, 11, 11, 3, 9, 10, 10, 12, 2, 12, 5, 10},
	{1, 1, 1, 1, 11, 9, 11, 8, 9, 9, 9, 8, 4, 10, 12, 6, 7, 5, 12, 6, 11, 5, 11, 8, 10, 11, 12, 10, 3, 10, 8, 10, 8, 8, 10, 6, 8, 9, 7, 12, 9, 7, 4, 11, 12, 3, 12, 9, 11, 7, 7, 4, 7, 6, 7, 10, 5, 8, 3, 10, 11, 4, 12, 9, 3, 12, 5, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 50, 200, 1000}, -- wild
	[ 2] = {0, 0, 0, 0, 0},       -- scatter (2, 3, 4 reels only)
	[ 3] = {0, 0, 40, 150, 750},  -- white
	[ 4] = {0, 0, 40, 150, 750},  -- black
	[ 5] = {0, 0, 20, 80, 400},   -- blue amulet
	[ 6] = {0, 0, 20, 80, 400},   -- red amulet
	[ 7] = {0, 0, 15, 60, 300},   -- ace
	[ 8] = {0, 0, 15, 60, 300},   -- king
	[ 9] = {0, 0, 10, 40, 200},   -- queen
	[10] = {0, 0, 10, 40, 200},   -- jack
	[11] = {0, 0, 5, 20, 100},    -- ten
	[12] = {0, 0, 5, 20, 100},    -- nine
}

-- 3. FREE SPINS TABLE
local scat_fs = 10

-- 4. CONFIGURATION
local sx, sy = 5, 4 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
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
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
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
		local q = 0.5 -- probability of getting or not getting free spins is equal
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
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
