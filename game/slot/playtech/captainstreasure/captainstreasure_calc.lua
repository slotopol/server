-- Playtech / Captain's Treasure
-- RTP calculation

-- 1. REEL STRIPS DATA
local REELS = {
	-- luacheck: push ignore 631
	{11, 9, 8, 10, 4, 9, 8, 11, 10, 8, 7, 4, 10, 7, 8, 6, 7, 9, 8, 7, 11, 5, 10, 7, 9, 11, 6, 5, 3, 8, 6, 11, 9, 5, 6, 10, 7, 11, 9, 10, 2},
	{9, 8, 2, 11, 5, 8, 11, 5, 8, 6, 10, 7, 4, 9, 10, 6, 7, 11, 8, 9, 11, 6, 9, 10, 11, 3, 9, 10, 8, 4, 11, 10, 8, 7, 6, 9, 10, 7, 6, 1, 7},
	{9, 10, 7, 8, 9, 3, 7, 5, 11, 10, 7, 11, 9, 10, 6, 11, 4, 10, 8, 1, 6, 7, 8, 5, 4, 11, 7, 8, 9, 10, 5, 11, 7, 10, 9, 11, 6, 8, 2, 9, 8, 6},
	{2, 9, 5, 10, 11, 8, 7, 9, 8, 4, 10, 7, 9, 1, 11, 8, 10, 11, 8, 6, 10, 8, 6, 7, 3, 11, 9, 10, 6, 7, 11, 9, 6, 5, 9, 6, 4, 7, 11, 8, 10},
	{11, 9, 4, 7, 10, 9, 8, 10, 7, 11, 5, 9, 8, 3, 7, 10, 11, 5, 9, 11, 10, 6, 11, 8, 6, 2, 9, 7, 8, 10, 6, 5, 10, 7, 8, 11, 6, 9, 8, 4, 7},
	-- luacheck: pop
}

-- 2. PAYTABLE FOR LINE WINS (indexed by symbol ID)
local PAYTABLE_LINE = {
	[ 1] = {0, 0, 0, 0, 0},         -- wild (2, 3, 4 reels only)
	[ 2] = {0, 0, 0, 0, 0},         -- scatter
	[ 3] = {2, 10, 100, 500, 5000}, -- sabers
	[ 4] = {0, 5, 50, 250, 2500},   -- map
	[ 5] = {0, 3, 20, 100, 1000},   -- anchor
	[ 6] = {0, 0, 10, 30, 500},     -- ace
	[ 7] = {0, 0, 5, 25, 300},      -- king
	[ 8] = {0, 0, 5, 20, 200},      -- queen
	[ 9] = {0, 0, 5, 20, 200},      -- jack
	[10] = {0, 0, 5, 15, 100},      -- ten
	[11] = {0, 0, 5, 15, 100},      -- nine
}

-- 3. PAYTABLE FOR SCATTER WINS (for 1 selected line bet)
local PAYTABLE_SCAT = {0, 1, 5, 10, 100}
local scat_min = 2 -- minimum scatters to win

-- 4. CONFIGURATION
local sx, sy = 5, 3 -- screen width & height
local wild, scat = 1, 2 -- wild & scatter symbol IDs

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

		for symbol_id, pays in pairs(PAYTABLE_LINE) do
			if symbol_id ~= wild and symbol_id ~= scat then
				local s = symbol_counts(symbol_id)
				local c = {}
				for i = 1, sx do c[i] = s[i] + w[i] end

				local function get_ev_for_direction(is_left_to_right)
					local dir_ev = 0
					for n = 1, 5 do
						-- If it's a 5-of-a-kind and we're counting from right to left,
						-- we skip it to avoid counting the same line twice.
						if pays[n] > 0 and (n < 5 or is_left_to_right) then
							local total_combs = 1
							local no_wild_combs = 1

							for i = 1, 5 do
								-- Determine the reel index depending on the direction
								local idx = is_left_to_right and i or (6 - i)
								if i <= n then
									total_combs = total_combs * c[idx]
									no_wild_combs = no_wild_combs * s[idx]
								elseif i == n + 1 then
									total_combs = total_combs * (lens[idx] - c[idx])
									no_wild_combs = no_wild_combs * (lens[idx] - c[idx])
								else
									total_combs = total_combs * lens[idx]
									no_wild_combs = no_wild_combs * lens[idx]
								end
							end

							local with_wild = total_combs - no_wild_combs
							dir_ev = dir_ev + (no_wild_combs * pays[n]) + (with_wild * pays[n] * 2)
						end
					end
					return dir_ev
				end

				ev_sum = ev_sum + get_ev_for_direction(true)  -- Left to Right
				ev_sum = ev_sum + get_ev_for_direction(false) -- Right to Left
			end
		end

		return ev_sum
	end

	-- Function to calculate expected return from scatter wins
	local function calculate_scat_ev()
		local c = symbol_counts(scat)
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
