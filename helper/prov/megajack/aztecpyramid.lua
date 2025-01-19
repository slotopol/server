local path = arg[0]:match("(.*[/\\])")

local function printf(fmt, ...)
  print(string.format(fmt, ...))
end

local apm = {10, 12, 15, 20, 25, 51}
local att = 3 -- number of attempts

local s1 = 0
for i = 1,#apm do
  s1 = s1 + apm[i]
end

local s2 = 0
for i = 1,#apm do
  s2 = s2 + s1/apm[i]
end
printf("s1 = %d, s2 = %g", s1, s2)

print "probability of a single occurrence"

local apf = {} -- aztec pyramid frequency
for i = 1,#apm do
  apf[i] = s1/apm[i]/s2
end

for i, f in pairs(apf) do
  printf("pyramid #%d freq = %.12f", i, f)
end

print "probability of getting in 4 attempts"

local p6 = 1-apf[6]
local app = { -- aztec pyramid probability
  apf[1]*p6^(att-1),
  apf[2]*p6^(att-1),
  apf[3]*p6^(att-1),
  apf[4]*p6^(att-1),
  apf[5]*p6^(att-1),
  1-p6^att,
}

for i, p in pairs(app) do
  printf("pyramid #%d freq = %.12f", i, p)
end

math.randomseed(os.time())
local function getpyramid(p)
  local s = 0
  for i, f in ipairs(apf) do
    s = s + f
    if p <= s then
      return i
    end
  end
  return 6
end

print "check up by Monte Carlo method"
local res = {0, 0, 0, 0, 0, 0}
local num = 1000000
for _ = 1,num do
  local pyr
  for _ = 1, att do
    pyr = getpyramid(math.random())
    if pyr == 6 then
      break
    end
  end
  res[pyr] = res[pyr]+1
end
printf("p[1]=%g, p[2]=%g, p[3]=%g, p[4]=%g, p[5]=%g, p[6]=%g",
  res[1]/num, res[2]/num, res[3]/num, res[4]/num, res[5]/num, res[6]/num)

-- for 4 attempts
-- 10000m:
-- p[1]=0.2350307611
-- p[2]=0.1958578232
-- p[3]=0.1566851137
-- p[4]=0.1175131771
-- p[5]=0.0940089364
-- p[6]=0.2009041885
-- 1000000m:
-- p[1]=0.235028112811
-- p[2]=0.19585703166
-- p[3]=0.156685147831
-- p[4]=0.117513980418
-- p[5]=0.094012023438
-- p[6]=0.200903703842

local Epyr = 0
for i, p in pairs(app) do
  Epyr = Epyr + p*apm[i]
end
printf("pyramid expectation Epyr = %.12g", Epyr)

local Eroom = dofile(path.."aztecroom.lua")
local Ebon = Epyr + Eroom*app[6]
printf("Ebon = Epyr + Eroom*app[6] = %g + %g * %g = %.12g", Epyr, Eroom, app[6], Ebon)
