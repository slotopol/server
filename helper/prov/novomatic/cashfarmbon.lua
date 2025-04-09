
local function nextset(a)
	local j = #a - 1
	while j ~= 0 and a[j] >= a[j+1] do
		j = j - 1
	end
	if j == 0 then
		return false
	end
	local k = #a
	while a[j] >= a[k] do
		k = k - 1
	end
	a[j], a[k] = a[k], a[j]
	local l, r = j + 1, #a
	while l < r do
		a[l], a[r] = a[r], a[l]
		l, r = l + 1, r - 1
	end
	return true
end

local function sequence(a)
	return coroutine.wrap(function ()
		local i = 1
		repeat
			coroutine.yield(i, a)
			i = i + 1
		until not nextset(a)
	end)
end

-- average value of array
local function avr(r)
	local s = 0
	for _, v in pairs(r) do
	  s = s + v
	end
	return s / #r
end

local function shuffle(t)
	for i = #t, 1, -1 do
		local j = math.random(i)
		t[i], t[j] = t[j], t[i]
	end
end

-- bonus rules:
-- level1: 6 cells, 3 pay arrow, reel1
-- level2: 5 cells, $, 1 pay arrow, 1 mult arrow, reel1
-- level3: 4 cells, $, 1 pay arrow, 1 mult arrow, reel1
-- level4: 3 cells, $, 1 mult arrow, reel2
-- level5: 2 cells, $, reel2

math.randomseed(os.time())

local reel1 = {2, 2, 3, 4, 5, 6, 8, 10, 14}
local reel4 = {14, 14, 20, 24}
local reel5 = {20, 24, 40, 50}

-- expectation pay on reel
local er1 = avr(reel1)
local er4 = avr(reel4)
local er5 = avr(reel5)
print("Ereel1="..er1)
print("Ereel4="..er4)
print("Ereel5="..er5)

local t1 = {er1, er1, er1, -er1, -er1, -er1} -- level 1
local t2 = {er1, er1, -er1, -1, 0} -- level 2
local t3 = {er1, -er1, -1, 0} -- level 3
local t4 = {er4, -1, 0} -- level 4
local t5 = {er5, 0} -- level 5

local function calcbon(bon)
	local p, m = 0, 1
	for _, lev in ipairs(bon) do
		for _, c in ipairs(lev) do
			if c == 0 then
				return p, m
			elseif c == 1 or c == -1 then
				m = m + 1
			elseif c < 0 then
				p = p - c
			else
				p = p + c
			end
			if c < 0 then
				break
			end
		end
	end
	return p, m
end

local function bruteforce()
	local b1, b2, b3, b4, b5 = {}, {}, {}, {}, {}
	local bon = {b1, b2, b3, b4, b5}
	local E = 0
	for _, lev1 in sequence{1, 2, 3, 4, 5, 6} do
		b1[1], b1[2], b1[3], b1[4], b1[5], b1[6] =
			t1[lev1[1]], t1[lev1[2]], t1[lev1[3]], t1[lev1[4]], t1[lev1[5]], t1[lev1[6]]
		for _, lev2 in sequence{1, 2, 3, 4, 5} do
			b2[1], b2[2], b2[3], b2[4], b2[5] =
				t2[lev2[1]], t2[lev2[2]], t2[lev2[3]], t2[lev2[4]], t2[lev2[5]]
			for _, lev3 in sequence{1, 2, 3, 4} do
				b3[1], b3[2], b3[3], b3[4] =
					t3[lev3[1]], t3[lev3[2]], t3[lev3[3]], t3[lev3[4]]
				for _, lev4 in sequence{1, 2, 3} do
					b4[1], b4[2], b4[3] =
						t4[lev4[1]], t4[lev4[2]], t4[lev4[3]]
					for _, lev5 in sequence{1, 2} do
						b5[1], b5[2] =
							t5[lev5[1]], t5[lev5[2]]
						local p, m = calcbon(bon)
						E = E + p*m
					end
				end
			end
		end
	end
	local n = 2 * 6 * 24 * 120 * 720 -- 24883200
	return E/n
end

local function montecarlo()
	local E = 0
	local n = 10000000
	local bon = {t1, t2, t3, t4, t5}
	for _ = 1, n do
		shuffle(bon[1])
		shuffle(bon[2])
		shuffle(bon[3])
		shuffle(bon[4])
		shuffle(bon[5])
		local p, m = calcbon(bon)
		E = E + p*m
	end
	return E/n
end

local E = bruteforce()
local _ = montecarlo

-- Ebon:	50
print("Ebon:", E)
