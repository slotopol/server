local reelmt = {}
reelmt.__index = function(t, i)
	return rawget(t, (i - 1) % rawlen(t) + 1)
end
reelmt.__newindex = function(t, i, v)
	rawset(t, (i - 1) % rawlen(t) + 1, v)
end

function makereel(symset, neighbours)
	-- make not-shuffled reel
	local reel = {}
	for sym, n in pairs(symset) do
		for _ = 1, n do
			table.insert(reel, sym)
		end
	end
	setmetatable(reel, reelmt)

	-- shuffle it
	shuffle(reel)

	-- correct it
	local iter = correctreel(reel, neighbours)
	return reel, iter
end

function makereelhot(symset, sy, scat, chunklen, strict)
	-- make set of chunks
	local chunks = {}
	for sym, n in pairs(symset) do
		while n > 0 do
			local m
			if scat[sym] then
				m = chunklen[sym] or 1
			else
				m = math.min(n, chunklen[sym] or sy)
			end
			n = n - m
			local c = {sym=sym, n=m}
			chunks[#chunks+1] = c
		end
	end

	-- shuffle until reel become correct
	local iter = 0
	repeat
		shuffle(chunks)
		local ok, n = true, 0
		local last = 0
		for i = 1, #chunks do
			local c1, c2 = chunks[i], chunks[i+1]
			if strict and c1.sym == last then
				if not c2 or c2.sym == last then
					ok = false
					break
				else
					c1, c2 = c2, c1
					chunks[i], chunks[i+1] = c1, c2
				end
			end
			if scat[c1.sym] then
				if n < sy then
					if not c2 or scat[c2.sym] or (strict and c2.sym == last) then
						ok = false
						break
					else
						c1, c2 = c2, c1
						chunks[i], chunks[i+1] = c1, c2
					end
				else
					n = 0
				end
			else
				n = n + c1.n
			end
			last = c1.sym
		end
		iter = iter + 1
	until ok or iter >= 1000

	-- glue chunks to single reel
	local reel = {}
	for _, c in pairs(chunks) do
		for _ = 1, c.n do
			reel[#reel+1] = c.sym
		end
	end
	return reel, iter
end

function shuffle(t)
	for i = #t, 1, -1 do
		local j = math.random(i)
		t[i], t[j] = t[j], t[i]
	end
end

local maxiter = 10000

function correctreel(reel, neighbours)
	local iter = 0
	while true do
		local n = 0
		for i = 1, #reel do
			local symi = reel[i]
			local b
			b = neighbours[symi][reel[i - 3]]
			if b >= 3 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-3 then j = j+7 end
				reel[i - 3], reel[j] = reel[j], reel[i - 3]
				n = n + 1
			end
			b = neighbours[symi][reel[i - 2]]
			if b >= 2 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-2 then j = j+5 end
				reel[i - 2], reel[j] = reel[j], reel[i - 2]
				n = n + 1
			end
			b = neighbours[symi][reel[i - 1]]
			if b >= 1 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-1 then j = j+3 end
				reel[i - 1], reel[j] = reel[j], reel[i - 1]
				n = n + 1
			end
			b = neighbours[symi][reel[i + 1]]
			if b >= 1 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-1 then j = j+3 end
				reel[i + 1], reel[j] = reel[j], reel[i + 1]
				n = n + 1
			end
			b = neighbours[symi][reel[i + 2]]
			if b >= 2 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-2 then j = j+5 end
				reel[i + 2], reel[j] = reel[j], reel[i + 2]
				n = n + 1
			end
			b = neighbours[symi][reel[i + 3]]
			if b >= 3 then
				local j = math.random(#reel - b * 2 - 1)
				if j >= i-3 then j = j+7 end
				reel[i + 3], reel[j] = reel[j], reel[i + 3]
				n = n + 1
			end
		end
		iter = iter + 1
		if n == 0 then
			break
		end
		if iter >= maxiter then
			break
		end
	end
	return iter
end

function printreel(reel, iter)
	if iter and iter > 1 then
		if iter >= maxiter then
			print"too many neighbours shuffle iterations"
			return
		else
			print(iter.." iterations")
		end
	end
	print("- [" .. table.concat(reel, ", ") .. "] # "..#reel) -- for yaml-file
end
