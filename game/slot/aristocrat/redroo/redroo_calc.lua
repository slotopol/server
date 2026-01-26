-- Aristocrat / Redroo
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{5, 10, 11, 9, 13, 4, 8, 10, 11, 4, 12, 10, 11, 9, 12, 2, 13, 6, 8, 10, 12, 7, 13, 5, 11, 7, 8, 9, 5, 13, 8, 4, 9, 6, 13, 3, 12, 7, 11, 6, 10, 7, 9, 6, 8, 3, 12, 5, 13, 3, 10, 7, 9, 4, 12, 11, 6, 8},
	{7, 12, 6, 9, 10, 4, 12, 8, 10, 11, 5, 12, 2, 13, 5, 8, 4, 13, 7, 8, 4, 11, 6, 9, 3, 12, 7, 9, 3, 12, 13, 11, 10, 5, 13, 6, 9, 11, 1, 10, 6, 8, 4, 12, 6, 11, 10, 8, 13, 7, 11, 5, 9, 3, 13, 7, 8, 9, 10},
	{4, 12, 7, 10, 6, 9, 5, 8, 6, 13, 5, 9, 6, 10, 1, 9, 7, 12, 11, 8, 3, 12, 10, 8, 4, 13, 11, 5, 10, 4, 11, 9, 10, 6, 13, 11, 12, 7, 13, 5, 8, 2, 10, 6, 11, 3, 13, 12, 9, 11, 3, 12, 13, 7, 8, 4, 9, 7, 8},
	{13, 3, 8, 10, 6, 11, 13, 4, 12, 6, 13, 8, 11, 7, 10, 4, 9, 13, 6, 12, 11, 2, 9, 7, 11, 8, 10, 4, 9, 6, 12, 11, 5, 9, 13, 3, 10, 5, 11, 3, 10, 5, 13, 7, 8, 12, 9, 7, 8, 12, 6, 9, 4, 10, 7, 8, 1, 12, 5},
	{13, 5, 10, 4, 8, 6, 12, 5, 13, 9, 3, 12, 6, 13, 11, 8, 4, 10, 12, 9, 11, 3, 10, 8, 13, 5, 11, 7, 10, 2, 11, 7, 12, 9, 5, 8, 6, 10, 7, 9, 4, 12, 11, 8, 9, 13, 7, 12, 4, 11, 9, 6, 10, 3, 13, 7, 8, 6},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                     --  1 wild (2, 3, 4 reels only)
	[ 2] = {},                     --  2 scatter
	[ 3] = {0, 50, 150, 200, 250}, --  3 redroo
	[ 4] = {0, 20, 80, 150, 200},  --  4 girl
	[ 5] = {0, 20, 80, 150, 200},  --  5 boy
	[ 6] = {0, 10, 40, 100, 150},  --  6 dog
	[ 7] = {0, 10, 40, 100, 150},  --  7 parrot
	[ 8] = {0, 0, 10, 50, 140},    --  8 ace
	[ 9] = {0, 0, 10, 50, 140},    --  9 king
	[10] = {0, 0, 5, 40, 120},     -- 10 queen
	[11] = {0, 0, 5, 40, 120},     -- 11 jack
	[12] = {0, 0, 5, 20, 100},     -- 12 ten
	[13] = {0, 2, 5, 20, 100},     -- 13 nine
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 0, 80, 400, 800}
local scat_min = 2 -- minimum scatters to win

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT_REG = {0, 0, 8, 15, 20}
local FREESPIN_SCAT_BON = {0, 5, 8, 15, 20}

-- 5. CONFIGURATION
local sx, sy = 5, 4 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local cost = 50 -- cost of spin with bet=1
local line_min = 2 -- minimum line symbols to win
local mw = 2.5 -- average multiplier on wilds

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

	-- Function to calculate expected return by ways for all symbols
	local function calculate_ways_ev(free_spins)
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
								if s == wild and free_spins then
									n = n + mw
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
	local function calculate_scat_ev(free_spins)
		local c = counts[scat]
		local ev_sum, fs_sum, fs_num = 0, 0, 0
		local FREESPIN_SCAT = free_spins and FREESPIN_SCAT_BON or FREESPIN_SCAT_REG

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
					if FREESPIN_SCAT[scat_sum] > 0 then
						fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
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

	-- Execute calculation
	local rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	do
		local rtp_line = calculate_ways_ev(true) / reshuffles / cost * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(true)
		local rtp_scat = ev_sum / reshuffles / cost * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_fs = sq * rtp_sym
		print(string.format("*free games calculations*"))
		print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
		print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
		print(string.format("free games frequency: 1/%.5g", reshuffles/fs_num))
		print(string.format("RTP = sq*rtp(sym) = %.5g*%.5g = %.6f%%", sq, rtp_sym, rtp_fs))
	end
	local rtp_total
	do
		local rtp_line = calculate_ways_ev(false) / reshuffles / cost * 100
		local ev_sum, fs_sum, fs_num = calculate_scat_ev(false)
		local rtp_scat = ev_sum / reshuffles / cost * 100
		local rtp_sym = rtp_line + rtp_scat
		local q = fs_sum / reshuffles
		local sq = 1 / (1 - q)
		rtp_total = rtp_sym + q * rtp_fs
		print(string.format("*regular games calculations*"))
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

calculate(REELS)
