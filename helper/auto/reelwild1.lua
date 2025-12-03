
-- This script is for reels sets composition by reshuffles reels content.
-- Useful for games with reel wilds or some others cases where symbols
-- reshuffle at the reel gets new RTP. Implemented by sequential scanner
-- run for each reels set.

--- input data begin ---

-- path to slotopol executable file, place here others necessary flags
local slotpath = "slot_debug.exe"
-- provider / gamename
local gamename = "ctinteractive/hellshot7s"
-- path to reels generator script
local gamescript = "ct/hellshot7s.lua"
-- number of reels at videoslot
local reelnum = 5
-- number of generator iterations
local N = 100
-- RTP granulation, can be 0.5, 1.0, 2.0
local gran = 0.5

-- temporary yaml file to check up by scanner
local genfile = (os.getenv("TEMP") or os.getenv("TMP") or os.getenv("TMPDIR")).."/reelgen.yaml"
-- final yaml file name and path
local devfile = os.getenv("GOPATH").."/bin/reeldev.yaml"

--- input data end ---

autoscan = true
local scripts = arg[0]:match("^(.*helper[/%\\])")
local reelgen = dofile(scripts.."prov/"..gamescript)
assert(type(reelgen) == "function", "reels generator function 'reelgen' does not defined")

local keypool = {}

local function generate()
	-- make reels set
	local reels = {}
	for i = 1, reelnum do
		reels[i] = reelgen(i)
	end

	-- write temporary yaml-file
	local f, err = io.open(genfile, "w")
	if not f then
		error("cannot create generator file: "..err)
	end
	f:write("\n", gamename.."/reel\n\n---\n\n50:\n")
	for _, reel in ipairs(reels) do
		f:write("  - [" .. table.concat(reel, ", ") .. "] # "..rawlen(reel).."\n")
	end
	f:close()

	-- run scanner
	local cltpl = "%s -f=\"%s\" scan --noembed -g=\"%s@50\"" -- command line template
	local cl = string.format(cltpl, slotpath, genfile, gamename) -- command line parameters
	local h = io.popen(cl)
	if not h then
		error("cannot run command: "..cl)
	end
	local output = h:read("*a")
	h:close()

	reels.comment = assert(output:match("(reels lengths.*)$"), "calculation output does not received")
	reels.rtp = assert(
		tonumber(string.match(output, " = (%d+%.?%d-)%%%s$")),
		"result RTP does not found, comment is:\n"..reels.comment)

	return reels
end

-- run scanner N times
for stage = 1, N do
	local reels = generate()
	reels.diff = reels.rtp % gran
	local key = string.format("%.1f", reels.rtp - reels.diff)
	if not keypool[key] or keypool[key].diff > reels.diff then
		keypool[key] = reels
	end
	print(string.format("(%d/%d) RTP = %g%%", stage, N, reels.rtp))
end

-- make sorted table with granulated reels sets
local t, i = {}, 1
for _, reels in pairs(keypool) do
	t[i], i = reels, i + 1
end
table.sort(t, function(a, b) return a.rtp < b.rtp end)

-- write final yaml-file
local f, err = io.open(devfile, "w")
if not f then
	error("cannot create results file: "..err)
end
f:write("\n", gamename.."/reel\n\n---\n")
for _, reels in pairs(t) do
	f:write("\n", reels.comment:gsub("(.-)\n", "# %1\n")..reels.rtp..":\n")
	for _, reel in ipairs(reels) do
		f:write("  - [" .. table.concat(reel, ", ") .. "] # "..rawlen(reel).."\n")
	end
end
f:close()

print(string.format("%d entries in complete file, spent %g seconds", #t, os.clock()))

os.remove(genfile)
