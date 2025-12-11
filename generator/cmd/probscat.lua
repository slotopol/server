
-- Function to calculate the probability of getting between N1 and N2 scatters (inclusive) across 15 spins.
-- p_reels is a table/array of 5 probabilities (p1, p2, p3, p4, p5).
local function calculate_range_probability(p_reels, N1, N2)
	local num_spins = 15
	local max_scatters_total = num_spins * 5 -- Maximum possible is 75 scatters total

	-- Step 1: Calculate the probability distribution for a SINGLE spin (0 to 5 scatters)
	-- current_dist[m] = probability of getting exactly 'm' scatters at the current stage
	local current_dist = { [0] = 1.0 }

	for _, p in ipairs(p_reels) do
		local next_dist = {}
		for scatters, prob in pairs(current_dist) do
			-- Case A: Scatter does not hit (1-p probability)
			next_dist[scatters] = (next_dist[scatters] or 0.0) + prob * (1.0 - p)
			-- Case B: Scatter hits (p probability), count increases by 1
			next_dist[scatters + 1] = (next_dist[scatters + 1] or 0.0) + prob * p
		end
		current_dist = next_dist
	end

	-- P_m will be an array (indices 0 to 5) of probabilities for a single spin
	local P_m = {}
	for m = 0, 5 do
		P_m[m] = current_dist[m] or 0.0
	end

	-- Step 2: Calculate the probability distribution for 15 spins using dynamic programming (convolution)
	-- dp[k] = probability of getting exactly k total scatters after 'i' spins
	local dp = {}
	-- At the start (0 spins), 100% probability of 0 scatters
	dp[0] = 1.0

	for _ = 1, num_spins do
		local next_dp = {}
		for k = 0, max_scatters_total do -- k is total scatters so far
			if dp[k] and dp[k] > 0 then
				for m = 0, 5 do -- m is scatters on the current spin
					if k + m <= max_scatters_total then
						-- P(total k+m) += P(total k) * P(current m)
						next_dp[k + m] = (next_dp[k + m] or 0.0) + dp[k] * P_m[m]
					end
				end
			end
		end
		dp = next_dp
	end
	-- The dp array now holds probabilities for all possible outcomes (0 to 75 total scatters)

	-- Step 3: Sum probabilities for the range N1 to N2
	local total_probability = 0.0

	-- Ensure N1 and N2 are within valid bounds [0, max_scatters_total]
	local start_k = math.max(0, N1)
	local end_k = math.min(max_scatters_total, N2)

	for k = start_k, end_k do
		if dp[k] then
			total_probability = total_probability + dp[k]
		end
	end

	return total_probability
end

-- Example Usage:
-- Set your probabilities p1, p2, p3, p4, p5
local probabilities = {3/30, 3/34, 3/32, 3/33, 3/30}

-- Define the range [N1, N2]. E.g., for 11 or more, set N1=11, N2=75 (max possible)
local N1_target = 11
local N2_target = 75 -- Max possible scatters is 15 * 5 = 75

local result = calculate_range_probability(probabilities, N1_target, N2_target)

print(string.format("probability of getting %d to %d scatters in 15 spins: %.8f, freq: %g",
	N1_target, N2_target, result, 1/result))

-- Example for exactly 11 scatters: N1=11, N2=11
local result_exact = calculate_range_probability(probabilities, 11, 11)
print(string.format("probability of getting exactly 11 scatters: %.8f", result_exact))
