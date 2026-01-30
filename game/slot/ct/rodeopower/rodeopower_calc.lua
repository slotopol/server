-- CT Interactive / Rodeo Power
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS_BON = {
	-- luacheck: push ignore 631
	{7, 4, 11, 12, 9, 13, 8, 10, 4, 12, 9, 3, 5, 11, 12, 8, 5, 4, 6, 12, 3, 6, 13, 9, 10, 7, 5, 3, 7, 6, 13, 2, 10, 11, 13, 2, 8},
	{4, 7, 5, 13, 3, 12, 7, 3, 12, 4, 5, 11, 12, 13, 9, 10, 8, 6, 10, 9, 6, 1, 12, 8, 13, 11, 10, 7, 3, 9, 5, 4, 2, 8, 13, 11, 6},
	{12, 4, 3, 9, 5, 10, 4, 2, 12, 6, 13, 11, 12, 5, 3, 9, 11, 10, 9, 6, 3, 13, 6, 11, 12, 8, 5, 7, 8, 2, 7, 13, 8, 4, 13, 10, 7},
	{10, 12, 13, 7, 2, 3, 13, 9, 10, 11, 8, 9, 12, 5, 7, 6, 11, 5, 4, 13, 6, 12, 3, 4, 5, 9, 8, 10, 11, 7, 12, 3, 8, 1, 13, 4, 6},
	{8, 10, 2, 11, 9, 8, 3, 12, 11, 7, 4, 6, 5, 12, 10, 13, 11, 6, 3, 10, 8, 9, 4, 3, 13, 9, 5, 12, 4, 7, 13, 2, 5, 7, 6, 13, 12},
	-- luacheck: pop
}
local REELS_REG = {
	-- luacheck: push ignore 631
	{13, 10, 7, 4, 6, 12, 4, 3, 8, 13, 12, 11, 10, 5, 9, 11, 10, 7, 13, 8, 9, 11, 2, 8, 11, 6, 10, 12, 5, 7, 8, 6, 13, 9, 3, 12, 9},
	{12, 3, 7, 8, 13, 9, 4, 10, 13, 1, 11, 9, 12, 3, 10, 5, 7, 8, 6, 13, 11, 9, 7, 13, 8, 10, 12, 11, 6, 4, 11, 5, 12, 9, 10, 2, 8, 6},
	{13, 8, 7, 4, 9, 11, 12, 6, 9, 12, 4, 6, 5, 8, 11, 13, 10, 12, 13, 10, 8, 7, 2, 3, 8, 10, 11, 6, 10, 5, 9, 13, 3, 9, 12, 11, 7},
	{7, 11, 1, 5, 7, 10, 9, 11, 12, 13, 9, 12, 13, 7, 11, 4, 3, 9, 10, 8, 6, 13, 11, 5, 12, 10, 3, 8, 13, 2, 9, 7, 8, 10, 6, 8, 4, 12},
	{5, 13, 12, 8, 4, 13, 7, 8, 10, 11, 6, 5, 13, 9, 10, 13, 4, 6, 10, 11, 8, 7, 9, 12, 3, 6, 10, 11, 9, 8, 2, 12, 7, 3, 9, 11, 12},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                    --  1 wild (2, 4 reels only)
	[ 2] = {},                    --  2 scatter
	[ 3] = {0, 0, 50, 300, 1000}, --  3 shoe
	[ 4] = {0, 0, 35, 300, 500},  --  4 woman
	[ 5] = {0, 0, 25, 250, 400},  --  5 spurs
	[ 6] = {0, 0, 25, 250, 400},  --  6 belt
	[ 7] = {0, 0, 10, 20, 120},   --  7 saddle
	[ 8] = {0, 0, 10, 20, 120},   --  8 hat
	[ 9] = {0, 0, 10, 20, 120},   --  9 boots
	[10] = {0, 0, 5, 10, 100},    -- 10 ace
	[11] = {0, 0, 5, 10, 100},    -- 11 king
	[12] = {0, 0, 5, 10, 100},    -- 12 queen
	[13] = {0, 0, 5, 10, 100},    -- 13 jack
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 5, 20, 100}
local scat_min = 3 -- minimum scatters to win
local scat_fs = 15 -- number of free spins awarded

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local cost = 25 -- cost of spin with bet=1
local line_min = 3 -- minimum line symbols to win

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

	-- Function to calculate expected return by ways for all symbols
	local function calculate_ways_ev()
		local ev_sum = 0

		for sym_id, pays in pairs(PAYTABLE_LINE) do
			if sym_id ~= wild and #pays > 0 then
				local c, z = {}, {}
				for x, r in ipairs(reels) do
					local len = lens[x]
					c[x], z[x] = 0, 0
					for i = 1, len do
						local n = 0 -- count in window
						for h = 0, sy - 1 do
							local s = r[(i + h - 1) % len + 1]
							if s == sym_id or s == wild then
								if s == wild and x == 2 then
									n = n + 2
								elseif s == wild and x == 4 then
									n = n + 5
								else
									n = n + 1
								end
							end
						end
						if n > 0 then
							c[x] = c[x] + n -- ways at reel
						else
							z[x] = z[x] + 1 -- stops without any
						end
					end
				end

				for x = line_min, sx do
					local pay = pays[x]
					if pay > 0 then
						local ways = 1
						for i = 1, sx do
							if i <= x then
								ways = ways * c[i] -- ways on 1..x reels
							elseif i == x + 1 then
								ways = ways * z[i] -- stops without sym_id on reel[x+1]
							else
								ways = ways * lens[i] -- anything on remaining reels
							end
						end
						ev_sum = ev_sum + ways*pay
					end
				end
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
		local rtp_line = calculate_ways_ev() / reshuffles / cost * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / reshuffles / cost * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_fs = sq * rtp_sym
		print(string.format("*bonus reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%", sq, rtp_sym, rtp_fs))
	end
	local rtp_total
	do
		reels = reels_reg
		precalculate_reels()
		local rtp_line = calculate_ways_ev() / reshuffles / cost * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev()
		local rtp_scat = ev_sum / reshuffles / cost * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular reels calculations*"))
		print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	end
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS_REG, REELS_BON)
