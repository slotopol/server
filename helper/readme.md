
# Helper scripts

Scripts in this section helps to build sets of reels. Typical workflow of reels composition for any game is followed:

* Temporary insert into algorithm of the game reels loading from some external file. To do this, replace code `var ReelsMap = slot.ReadMap[*slot.Reels5x](reels)` to `var ReelsMap = slot.LoadMap[*slot.Reels5x]("/some/path/reeldev.yaml")` and compile it.

* Create new script at `helper/prov` folder for developed game. Code in it can call `makereelhot` script function for reels same as in `Sizzling Hot` game, or call `makereel` script function if symbols neighbours have the matter.

* You can run this script at `ZeroBrane Studio` or some other way to get the reel with scripted rules.

* Complete set of reels as it written at previous step, and place this set into external `reeldev.yaml` file with some percentage, 50 for example.

* Run the scanner for this reels by command same as `slot_win_x64.exe scan -g=provider/gamename -r=50`. If reels in set have a long length (60 symbols or more), you can quickly evaluate RTP of reels by Monte Carlo method. Run scanner with `--mc=20` parameter for example. If RTP is approximately suitable you can perform precise calculation.

* Insert onto final yaml-file calculated reels with calculated percentage.

See also: [issue #3](https://github.com/slotopol/server/issues/3)
