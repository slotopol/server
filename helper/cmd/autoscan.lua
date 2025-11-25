
--- input data begin ---

local slotpath = "slot_win_x64.exe" -- path to slotopol executable file
local gamename = "ctinteractive/hellshot7s" -- provider / gamename
local gamescript = "ct/hellshot7s.lua" -- path to reels generator script
local reelnum = 5 -- number of reels at videoslot
local N = 8 -- number of generator iterations
local gran = 0.5 -- RTP granulation, can be 0.5, 1.0, 2.0

-- temporary yaml file to check up by scanner
local genfile = (os.getenv("TEMP") or os.getenv("TMP") or os.getenv("TMPDIR")).."/reelgen.yaml"
-- final yaml file name and path
local devfile = os.getenv("GOPATH").."/bin/reeldev.yaml"
-- command line template
local cltpl = "%s -f=\"%s\" scan -g=\"%s\" -r=50"

--- input data end ---

autoscan = true
local scripts = arg[0]:match("^(.*[/%\\]helper[/%\\])")
dofile(scripts.."prov/"..gamescript)
assert(type(reelgen) == "function", "reels generator function 'reelgen' does not defined")

local devpool = {}

local function generate()
	-- make reels set
	local reels = {}
	for i = 1, reelnum do
		reels[i] = reelgen(i)
	end

	-- write temporary yaml-file
	local f = io.open(genfile, "w")
	f:write("\n", gamename.."/reel\n\n---\n\n50:\n")
	for _, reel in ipairs(reels) do
		f:write("  - [" .. table.concat(reel, ", ") .. "] # "..#reel.."\n")
	end
	f:close()

	-- run scanner
	local cl = string.format(cltpl, slotpath, genfile, gamename) -- command line
	local h = io.popen(cl)
	local output = h:read("*a")
	h:close()

	reels.comment = assert(output:match("(reels lengths.*)$"), "calculation output does not received")
	reels.rtp = assert(tonumber(string.match(output, " = (%d+%.%d+)%%%s$")), "result RTP does not found")

	return reels
end

-- run scanner N times
for i = 1, N do
	local reels = generate()
	reels.diff = reels.rtp % gran
	local key = string.format("%.1f", reels.rtp - reels.diff)
	if not devpool[key] or devpool[key].diff > reels.diff then
		devpool[key] = reels
	end
	print(string.format("%d/%d RTP = %g%%", i, N, reels.rtp))
end

-- make sorted table with granulated reels sets
local t, i = {}, 1
for _, reels in pairs(devpool) do
	t[i], i = reels, i + 1
end
table.sort(t, function(a, b) return a.rtp < b.rtp end)

-- write final yaml-file
local f = io.open(devfile, "w")
f:write("\n", gamename.."/reel\n\n---\n")
for _, reels in pairs(t) do
	f:write("\n", reels.comment:gsub("(.-)\n", "# %1\n")..reels.rtp..":\n")
	for _, reel in ipairs(reels) do
		f:write("  - [" .. table.concat(reel, ", ") .. "] # "..#reel.."\n")
	end
end
f:close()

print(#t.." entries in complete file")

os.remove(genfile)
