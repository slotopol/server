-- Novomatic / African Simba
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{9, 6, 11, 9, 10, 2, 12, 3, 10, 2, 8, 9, 11, 6, 9, 11, 4, 8, 5, 12, 10, 5, 7, 4, 8, 5, 7, 6, 11, 4, 10, 3, 12, 7, 6, 12, 3, 8, 5, 7},
	{12, 3, 7, 6, 9, 4, 12, 5, 10, 8, 5, 11, 6, 10, 5, 8, 3, 7, 11, 5, 8, 1, 7, 11, 3, 10, 6, 12, 7, 9, 4, 12, 9, 4, 11, 9, 6, 10, 8},
	{11, 4, 8, 10, 5, 12, 6, 10, 5, 11, 4, 7, 10, 6, 9, 3, 12, 2, 9, 12, 8, 6, 7, 3, 8, 9, 6, 7, 5, 10, 2, 11, 3, 7, 11, 4, 8, 12, 1, 9, 5},
	{10, 4, 7, 3, 11, 5, 8, 12, 6, 9, 3, 7, 1, 12, 3, 10, 9, 8, 6, 7, 12, 8, 6, 7, 5, 11, 8, 5, 11, 10, 9, 4, 10, 5, 11, 9, 4, 12, 6},
	{8, 3, 11, 6, 9, 4, 8, 12, 5, 7, 3, 10, 12, 5, 8, 12, 10, 4, 11, 9, 2, 7, 9, 12, 5, 10, 3, 11, 6, 7, 5, 8, 4, 9, 6, 10, 7, 6, 11, 2},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {},                     --  1 wild (2, 3, 4 reels only)
	[ 2] = {},                     --  2 scatter (1, 3, 5 reels only)
	[ 3] = {0, 0, 100, 500, 2500}, --  3 giraffe
	[ 4] = {0, 0, 50, 150, 750},   --  4 buffalo
	[ 5] = {0, 0, 25, 75, 250},    --  5 lemur
	[ 6] = {0, 0, 25, 75, 250},    --  6 flamingo
	[ 7] = {0, 0, 10, 25, 125},    --  7 ace
	[ 8] = {0, 0, 10, 25, 125},    --  8 king
	[ 9] = {0, 0, 10, 25, 125},    --  9 queen
	[10] = {0, 0, 5, 20, 100},     -- 10 jack
	[11] = {0, 0, 5, 20, 100},     -- 11 ten
	[12] = {0, 0, 5, 20, 100},     -- 12 nine
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local scat_fs = 12 -- number of free spins awarded
local scat_min = 3 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- grid width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs
local cost = 25 -- cost of spin with bet=1
local line_min = 3 -- minimum line symbols to win
local mfs = 3 -- multiplier on free spins

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
								n = n + 1
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
	local rtp_line = calculate_ways_ev() / reshuffles / cost * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles / cost * 100
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = mfs * sq * rtp_sym
	local rtp_total = rtp_sym + q * rtp_fs
	print(string.format("reels lengths [%s], total reshuffles %d", table.concat(lens, ", "), reshuffles))
	print(string.format("symbols: %.5g(lined) + %.5g(scatter) = %.6f%%", rtp_line, rtp_scat, rtp_sym))
	print(string.format("free spins %d, q = %.5g, sq = 1/(1-q) = %.6f", fs_sum, q, sq))
	print(string.format("free games hit rate: 1/%.5g", reshuffles/fs_num))
	print(string.format("RTP = %.5g(sym) + %.5g*%.5g(fg) = %.6f%%", rtp_sym, q, rtp_fs, rtp_total))
	return rtp_total
end

if autoscan then
	return calculate
end

calculate(REELS)
