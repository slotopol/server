
local function printf(fmt, ...)
	print(string.format(fmt, ...))
end

local room = {}

room[1] = {
	{ 10, 1},
	{ 11, 1},
	{ 12, 2},
	{ 13, 2},
	{ 14, 3},
	{ 15, 3},
	{ 16, 4},
	{ 10, 1},
	{ 11, 1},
	{ 12, 2},
	{ 13, 2},
	{ 14, 3},
	{ 15, 3},
	{ 10, 1},
	{ 11, 1},
	{ 12, 2},
	{ 13, 2},
	{ 14, 3},
	{ 10, 1},
	{ 11, 1},
	{ 12, 2},
	{ 13, 2},
	{  0,14},
}

room[2] = {
	{ 15, 5},
	{ 16, 5},
	{ 17, 6},
	{ 18, 6},
	{ 19, 7},
	{ 20, 8},
	{ 21, 8},
	{ 15, 5},
	{ 16, 5},
	{ 17, 6},
	{ 18, 6},
	{ 19, 7},
	{ 20, 8},
	{ 15, 5},
	{ 16, 5},
	{ 17, 6},
	{ 18, 6},
	{ 19, 7},
	{ 15, 5},
	{ 16, 5},
	{ 17, 6},
	{ 18, 6},
	{  0,14},
}

room[3] = {
	{ 20, 9},
	{ 21, 9},
	{ 22,10},
	{ 23,10},
	{ 24,11},
	{ 25,11},
	{ 26,12},
	{ 20, 9},
	{ 21, 9},
	{ 22,10},
	{ 23,10},
	{ 24,11},
	{ 25,11},
	{ 20, 9},
	{ 21, 9},
	{ 22,10},
	{ 23,10},
	{ 24,11},
	{ 20, 9},
	{ 21, 9},
	{ 22,10},
	{ 23,10},
	{  0,14},
}

room[4] = {
	{ 50,15},
	{ 50,15},
	{100,16},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,18},
}

room[5] = {
	{250,17},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
	{  0,13},
}

local sum1 = 0
for _, cell in pairs(room[1]) do
	if cell[2] ~= 14 then
		sum1 = sum1 + cell[1]
	end
end
sum1 = sum1 + 4*sum1/(#room[1]-1)
local Erow1 = sum1/#room[1]
printf("Erow1 = %.12g", Erow1)

local sum2 = 0
for _, cell in pairs(room[2]) do
	if cell[2] ~= 14 then
		sum2 = sum2 + cell[1]
	end
end
sum2 = sum2 + 4*sum2/(#room[2]-1)
local Erow2 = sum2/#room[2]
printf("Erow2 = %.12g", Erow2)

local sum3 = 0
for _, cell in pairs(room[3]) do
	if cell[2] ~= 14 then
		sum3 = sum3 + cell[1]
	end
end
sum3 = sum3 + 4*sum3/(#room[3]-1)
local Erow3 = sum3/#room[3]
printf("Erow3 = %.12g", Erow3)

local Erow4 = (50+50+100+0)/4
local p4 = 4/#room[4]
printf("Erow4*p4 = %.12g * %.12g = %.12g", Erow4, p4, Erow4*p4)

local Erow5 = 250
local p5 = 1/#room[5]
printf("Erow5*p5 = %.12g * %.12g = %.12g", Erow5, p5, Erow5*p5)

local Eroom = Erow1 + Erow2 + Erow3 + Erow4*p4 + Erow5*p5*p4
printf("Eroom = %.12g", Eroom)

return Eroom
