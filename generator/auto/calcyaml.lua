
-- This script helps to read generated YAML-file with reels data
-- and calculate pointed reels by script with RTP calculation function.

--- input data begin ---

-- data router ID for bonus reels if it has
local bonid = nil
-- data router ID for regular reels map
local regid = "novomatic/sizzlinghot/reel"
-- full path to YAML file with reels data
local yamlfile = "reeldev.yaml"
-- full path to calculation script
local calcscript = "game/slot/novomatic/sizzlinghot/sizzlinghot_calc.lua"
-- percent key pointed to regular reels to calculate
local rtplookup = 50

--- input data end ---

local function read_reels(iter)
	local reels_reg, reels_bon
	-- Skip comments and empty lines and return line with some data
	local bom = string.char(239, 187, 191)
	local function nextline(skip)
		repeat
			local line = iter()
			assert(line or skip, "unexpected end of file")
			if not line then
				return
			end
			if line:sub(1, 3) == bom then -- skip BOM
				line = line:sub(4)
			end
			if not line:match("^%s*#") and not line:match("^%s*$") then
				return line
			end
		until false
	end
	local function read_yaml_separator()
		assert(nextline() == "---", "expected yaml data separator ---")
	end
	-- Read plain reels data
	local function reels_lines()
		local t = {}
		for i = 1, 5 do
			t[i] = assert(nextline():match("- %[(.-)%]"))
		end
		return t
	end
	-- Convert to applicable reels
	local function convert_reels(rlines)
		local reels = {}
		for i, rline in ipairs(rlines) do
			local reel = {}
			for sym in rline:gmatch("%d+") do
				reel[#reel+1] = tonumber(sym)
			end
			reels[i] = reel
		end
		return reels
	end

	-- Read bonus reels if it has
	if bonid then
		assert(nextline() == bonid, "expected reels data ID "..bonid)
		read_yaml_separator()
		local rlines = reels_lines()
		reels_bon = convert_reels(rlines)
		read_yaml_separator()
	end

	-- Read regular reels map
	assert(nextline() == regid, "expected reels data ID "..regid)
	read_yaml_separator()
	local reelsmap = {}
	repeat
		local s = nextline(true)
		if not s then
			break
		end
		local n = assert(string.match(s, "^(%d+%.?%d*):$"), "rtp key does not found")
		local rtp = tonumber(n)
		reelsmap[rtp or 0] = reels_lines()
	until false

	-- Find closest reels to given look up key
	local rtp = -1000
	local rlines
	for p, v in pairs(reelsmap) do
		if math.abs(rtplookup-p) < math.abs(rtplookup-rtp) then
			rtp, rlines = p, v
		end
	end
	reels_reg = convert_reels(rlines)
	return reels_reg, reels_bon
end

local f, err = io.open(yamlfile, "r")
if not f then
	error("cannot open yaml file: "..err)
end
local iter = f:lines()
local reels_reg, reels_bon = read_reels(iter)
f:close()

autoscan = true
local calculate = dofile(calcscript)
calculate(reels_reg, reels_bon)
