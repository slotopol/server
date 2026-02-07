
# Generator scripts

Scripts in this section helps to build sets of reels. Typical workflow of reels composition for any game is followed:

* Create new script at `generator/prov` folder for developed game. Code in it can call `makereelhot` script function for reels same as in `Sizzling Hot` game, or call `makereel` script function if symbols neighbours have the matter.

* You can run this script at `ZeroBrane Studio` or some other way to get the reel with scripted rules. Change numbers of symbols at the table of lua-script to get reels with other percentage.

* Complete set of reels as it written at previous step, and place this set into external `/some/path/reeldev.yaml` file with some key percentage, 50 for example.

* Run the scanner for this reels by command same as `slot_win_x64.exe -f=/some/path/reeldev.yaml scan --noembed -g=provider/gamename -r=50`. If reels in the set have a long length (60 symbols or more), you can quickly evaluate RTP of reels by Monte Carlo method. Run scanner with `--mc=20` parameter for example. If RTP is approximately suitable you can perform precise calculation.

* Insert onto final yaml-file calculated reels with calculated percentage.

The content of `reeldev.yaml` file to develop new reels might look something like this:

```yaml

megajack/slotopol/reel

---

50:
  - [13, 1, 5, 12, 13, 11, 12, 11, 13, 8, 2, 12, 13, 3, 4, 6, 13, 2, 5, 10, 13, 9, 7, 8, 13, 10, 7, 9, 13, 3, 4, 6]
  - [9, 5, 10, 13, 9, 6, 3, 4, 13, 2, 12, 8, 12, 13, 11, 12, 11, 13, 5, 7, 10, 6, 3, 4, 13, 2, 12, 8, 13, 7, 1, 12]
  - [12, 13, 11, 12, 11, 13, 5, 10, 9, 7, 1, 12, 13, 3, 8, 6, 12, 13, 8, 4, 12, 2, 5, 10, 13, 7, 2, 13, 6, 3, 4, 9]
  - [12, 1, 2, 13, 6, 5, 12, 4, 8, 12, 13, 3, 10, 9, 7, 13, 11, 11, 11, 11, 13, 5, 12, 9, 8, 6, 13, 3, 10, 2, 7, 4]
  - [13, 11, 13, 12, 6, 4, 12, 3, 2, 5, 12, 10, 7, 12, 8, 1, 9, 12, 8, 9, 12, 4, 3, 12, 2, 5, 12, 10, 7, 13, 12, 6]

```

## Automation

In some cases, it's possible to automate the generation of reel sets. To do this, the Lua script must contain the `reelgen` function, which returns the reel for the specified number. Then, you can use one of the scripts in `generator/auto`.
