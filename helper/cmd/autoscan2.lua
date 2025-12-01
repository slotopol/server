
-- This script is for reels sets composition by reshuffles reels content.
-- Useful for games with reel wilds or some others cases where symbols
-- reshuffle at the reel gets new RTP. Implemented by single scanner call
-- for all reels sets, working in multistage mode.

--- input data begin ---

local slotpath = "slot_debug.exe" -- path to slotopol executable file
local gamename = "ctinteractive/hellshot7s" -- provider / gamename
local gamescript = "ct/hellshot7s.lua" -- path to reels generator script
local reelnum = 5 -- number of reels at videoslot
local N = 100 -- number of generator iterations
local gran = 0.5 -- RTP granulation, can be 0.5, 1.0, 2.0

-- temporary yaml file to check up by scanner
local genfile = (os.getenv("TEMP") or os.getenv("TMP") or os.getenv("TMPDIR")).."/reelgen.yaml"
-- final yaml file name and path
local devfile = os.getenv("GOPATH").."/bin/reeldev.yaml"

--- input data end ---

autoscan = true
local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
local reelgen = dofile(scripts.."prov/"..gamescript)
assert(type(reelgen) == "function", "reels generator function 'reelgen' does not defined")

local devpool = {}
local keypool = {}
local ridt = {} -- reels sets IDs table

-- write temporary yaml-file
do
	local f, err = io.open(genfile, "w")
	if not f then
		error("cannot create generator file: "..err)
	end
	f:write("\n", gamename.."/reel\n\n---\n")
	for stage = 1, N do
		ridt[stage] = string.format("-g=\"%s@%d\"", gamename, stage)
		-- make reels set
		local reels = {}
		for i = 1, reelnum do
			reels[i] = reelgen(i)
		end
		devpool[stage] = reels
		-- append reels to file
		f:write("\n", stage, ":\n")
		for _, reel in ipairs(reels) do
			f:write("  - [" .. table.concat(reel, ", ") .. "] # "..rawlen(reel).."\n")
		end
	end
	f:close()
end

-- run scanner
local cltpl = "%s -f=\"%s\" scan --noembed --lstage --vstage %s" -- command line template
local cl = string.format(cltpl, slotpath, genfile, table.concat(ridt, " ")) -- command line parameters
local h = io.popen(cl)
if not h then
	error("cannot run command: "..cl:sub(1, 100))
end
local output = h:read("*a")
h:close()

print(output)

-- extract comments and RTP from output for each reels sets
local pattern = "%((%d+)/"..N.."%) scan .-(reels lengths.-RTP = [^\n]-(%d+%.?%d-)%%%s)"
for stage, comment, rtp in output:gmatch(pattern) do
	stage, rtp = tonumber(stage), tonumber(rtp)
	local reels = assert(devpool[stage], "reels of stage "..stage.." does not found")
	reels.comment = comment
	reels.rtp = rtp
	reels.diff = reels.rtp % gran
	local key = string.format("%.1f", reels.rtp - reels.diff)
	if not keypool[key] or keypool[key].diff > reels.diff then
		keypool[key] = reels
	end
end

-- make sorted table with granulated reels sets
local t, i = {}, 1
for _, reels in pairs(keypool) do
	t[i], i = reels, i + 1
end
table.sort(t, function(a, b) return a.rtp < b.rtp end)

-- write final yaml-file
do
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
end

print(#t.." entries in complete file")

os.remove(genfile)
