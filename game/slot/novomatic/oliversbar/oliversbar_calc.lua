-- Novomatic / Oliver's Bar
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{10, 5, 12, 11, 7, 6, 11, 9, 3, 12, 7, 10, 3, 12, 8, 2, 11, 10, 12, 9, 5, 8, 9, 10, 6, 11, 4, 7, 5, 9, 6, 8, 2, 11, 8, 10, 9, 4, 13},
	{1, 1, 1, 5, 11, 4, 10, 11, 3, 12, 13, 7, 10, 11, 5, 10, 12, 9, 6, 10, 8, 6, 7, 11, 9, 8, 4, 10, 6, 9, 12, 2, 9, 3, 11, 12, 5, 9, 7, 2, 8},
	{1, 1, 1, 5, 9, 12, 11, 7, 9, 2, 11, 3, 8, 7, 2, 9, 12, 11, 6, 10, 13, 4, 12, 3, 8, 12, 6, 10, 9, 5, 11, 6, 10, 8, 5, 11, 8, 10, 4, 9, 10, 7},
	{1, 1, 1, 2, 10, 3, 11, 4, 10, 5, 12, 7, 6, 9, 11, 12, 9, 10, 6, 8, 11, 12, 7, 4, 11, 6, 10, 9, 2, 8, 5, 13, 12, 9, 5, 10, 8, 7, 9, 8, 3, 11},
	{3, 12, 5, 8, 12, 9, 7, 4, 12, 10, 3, 12, 8, 11, 5, 7, 11, 6, 9, 10, 11, 2, 10, 4, 11, 9, 8, 6, 7, 2, 10, 6, 9, 10, 11, 5, 9, 8, 13},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 0, 0, 0},        -- wild (2, 3, 4 reels only)
	[ 2] = {0, 5, 100, 500, 5000}, -- Oliver
	[ 3] = {0, 0, 25, 200, 1000},  -- friends
	[ 4] = {0, 0, 25, 200, 1000},  -- couple
	[ 5] = {0, 0, 15, 100, 500},   -- sweet-stuffs
	[ 6] = {0, 0, 15, 100, 500},   -- cocktails
	[ 7] = {0, 0, 10, 50, 200},    -- flower
	[ 8] = {0, 0, 10, 50, 200},    -- lime
	[ 9] = {0, 0, 5, 25, 100},     -- olives
	[10] = {0, 0, 5, 25, 100},     -- strawberries
	[11] = {0, 0, 5, 25, 100},     -- oranges
	[12] = {0, 2, 5, 25, 100},     -- cherry
	[13] = {0, 0, 0, 0, 0},        -- scatter
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 2, 5, 20, 500}
local scat_min = 2 -- minimum scatters to win

-- 4. FREE SPINS TABLE
local FREESPIN_SCAT = {0, 0, 20, 20, 20}

-- 5. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 13 -- wild & scatter symbol IDs

-- Performs full RTP calculation for given reels
local function calculate(reels)
	assert(#reels == sx, "unexpected number of reels")
	-- Get number of total reshuffles and lengths of each reel.
	local reshuffles, lens = 1, {}
	for i, r in ipairs(reels) do
		reshuffles = reshuffles * #r
		lens[i] = #r
	end

	-- Function to count symbol occurrences on each reel
	local function symbol_counts(symbol_id)
		local counts = {}
		for i, r in ipairs(reels) do
			local count = 0
			for _, sym in ipairs(r) do
				if sym == symbol_id then
					count = count + 1
				end
			end
			counts[i] = count
		end
		return counts
	end

	-- Function to calculate expected return from line wins for all symbols
	local function calculate_line_ev()
		local ev_sum = 0
		local w = symbol_counts(wild)

		-- Iterate through all symbols that pay on lines
		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			local s = symbol_counts(symbol_id)
			local function get_comb_ev(n, payout)
				if payout <= 0 then return 0 end
				local comb_ev = payout
				for i = 1, sx do
					if i <= n then
						comb_ev = comb_ev * (s[i] + w[i])
					elseif i == n + 1 then
						comb_ev = comb_ev * (lens[i] - (s[i] + w[i]))
					else
						comb_ev = comb_ev * lens[i]
					end
				end
				return comb_ev
			end
			for n = 2, sx do
				ev_sum = ev_sum + get_comb_ev(n, pays[n])
			end
		end

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev()
		local c = symbol_counts(scat)
		local ev_sum, fs_sum, fs_num = 0, 0, 0

		-- Using an recursive approach to sum combinations for exactly N scatters
		local function find_scatter_combs(reel_index, scat_sum, current_comb)
			if reel_index > sx then
				if scat_sum >= scat_min then
					ev_sum = ev_sum + current_comb * PAYTABLE_SCAT[scat_sum]
					fs_sum = fs_sum + current_comb * FREESPIN_SCAT[scat_sum]
					if FREESPIN_SCAT[scat_sum] > 0 then
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
	local rtp_line = calculate_line_ev() / reshuffles * 100
	local ev_sum, fs_sum, fs_num = calculate_scat_ev()
	local rtp_scat = ev_sum / reshuffles * 100
	local rtp_sym = rtp_line + rtp_scat
	local q = fs_sum / reshuffles
	local sq = 1 / (1 - q)
	local rtp_fs = 4 * sq * rtp_sym
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
