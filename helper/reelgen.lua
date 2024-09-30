local reelmt = {}
reelmt.__index = function(t, i)
	return rawget(t, (i - 1) % rawlen(t) + 1)
end
reelmt.__newindex = function(t, i, v)
	rawset(t, (i - 1) % rawlen(t) + 1, v)
end

function MakeReel(symset)
	local reel = {}
	for sym, n in ipairs(symset) do
		for _ = 1, n do
			table.insert(reel, sym)
		end
	end
	setmetatable(reel, reelmt)
	return reel
end

function ShuffleReel(reel)
	for i = #reel, 1, -1 do
		local j = math.random(i)
		reel[i], reel[j] = reel[j], reel[i]
	end
end

function CorrectReel(reel, neighbours)
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
		if iter >= 1000 then
			break
		end
	end
	return iter
end

function RrintReel(reel, iter)
	if iter > 1 then
		if iter >= 1000 then
			print"too many neighbours shuffle iterations"
			return
		else
			print(iter.." iterations")
		end
	end
	io.write("{")
	for i = 1, #reel do
		io.write(reel[i] .. ", ")
	end
	io.write("}\n")
end
