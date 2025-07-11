import {
    Sym,
    Pos,
    Wins,
    SlotGame,
    Slotx,
    Screen5x3,
    Reels5x,
    Line,
    BetLinesNetEnt5x3,
    ErrNoFeature,
    WinItem,
    FindClosest,
    ReadMap, // Assuming this would be used to load ReelsMap if not pre-assigned
    NewLinex,
    NewWins
} from "../../slot";

// The content of fruitshop_reel.yaml would be loaded into this map.
// For this conversion, we'll define its type and assume it's populated elsewhere.
// Example structure: { "96.71": new Reels5x(...) }
export let ReelsMap: { [key: string]: Reels5x } = {}; // To be populated, e.g. by an external loader or direct assignment

// Lined payment.
export const LinePay: number[][] = [
    [], // 0: empty, symbol 1 is wild
    [0, 5, 25, 300, 2000], // 1: Symbol 2 (Cherry)
    [0, 0, 25, 150, 1000], // 2: Symbol 3 (Plum)
    [0, 0, 20, 125, 750],  // 3: Symbol 4 (Lemon)
    [0, 0, 20, 100, 500],  // 4: Symbol 5 (Orange)
    [0, 0, 15, 75, 200],   // 5: Symbol 6 (Melon)
    [0, 0, 15, 50, 150],   // 6: Symbol 7 (Ace)
    [0, 0, 10, 25, 100],   // 7: Symbol 8 (King)
    [0, 0, 5, 20, 75],     // 8: Symbol 9 (Queen)
    [0, 0, 5, 15, 60],     // 9: Symbol 10 (Jack)
    [0, 0, 5, 10, 50],     // 10: Symbol 11 (Ten)
];

// Line freespins table on regular games
export const LineFreespinReg: number[][] = [
    [],              // 0: empty
    [0, 1, 1, 2, 5], // 1: Symbol 2 (Cherry)
    [0, 0, 1, 2, 5], // 2: Symbol 3 (Plum)
    [0, 0, 1, 2, 5], // 3: Symbol 4 (Lemon)
    [0, 0, 1, 2, 5], // 4: Symbol 5 (Orange)
    [0, 0, 1, 2, 5], // 5: Symbol 6 (Melon)
    [0, 0, 0, 0, 0], // 6: Symbol 7 (Ace)
    [0, 0, 0, 0, 0], // 7: Symbol 8 (King)
    [0, 0, 0, 0, 0], // 8: Symbol 9 (Queen)
    [0, 0, 0, 0, 0], // 9: Symbol 10 (Jack)
    [0, 0, 0, 0, 0], // 10: Symbol 11 (Ten)
];

// Line freespins table on bonus games
export const LineFreespinBon: number[][] = [
    [],              // 0: empty
    [0, 1, 1, 2, 5], // 1: Symbol 2 (Cherry)
    [0, 0, 1, 2, 5], // 2: Symbol 3 (Plum)
    [0, 0, 1, 2, 5], // 3: Symbol 4 (Lemon)
    [0, 0, 1, 2, 5], // 4: Symbol 5 (Orange)
    [0, 0, 1, 2, 5], // 5: Symbol 6 (Melon)
    [0, 0, 1, 2, 5], // 6: Symbol 7 (Ace)
    [0, 0, 1, 2, 5], // 7: Symbol 8 (King)
    [0, 0, 1, 2, 5], // 8: Symbol 9 (Queen)
    [0, 0, 1, 2, 5], // 9: Symbol 10 (Jack)
    [0, 0, 1, 2, 5], // 10: Symbol 11 (Ten)
];

// Bet lines
export const BetLines: Line[] = BetLinesNetEnt5x3.slice(0, 15);

const WILD_SYMBOL: Sym = 1;

export class Game extends Screen5x3 implements SlotGame {
    Slotx: Slotx;

    constructor() {
        super(); // Initializes Screen5x3 base (this.Screen)
        this.Slotx = new Slotx(BetLines.length, 1); // Sel, Bet
    }

    // Slotx getters/setters - delegating to composed Slotx object
    get Sel(): number { return this.Slotx.Sel; }
    set Sel(value: number) { this.Slotx.Sel = value; }
    get Bet(): number { return this.Slotx.Bet; }
    set Bet(value: number) { this.Slotx.Bet = value; }
    get Gain(): number | undefined { return this.Slotx.Gain; }
    set Gain(value: number | undefined) { this.Slotx.Gain = value; }
    get FSN(): number | undefined { return this.Slotx.FSN; }
    set FSN(value: number | undefined) { this.Slotx.FSN = value; }
    get FSR(): number | undefined { return this.Slotx.FSR; }
    set FSR(value: number | undefined) { this.Slotx.FSR = value; }

    Cost(): { cost: number, isJp: boolean } { return this.Slotx.Cost(); }
    Free(): boolean { return this.Slotx.Free(); }
    Spawn(wins: Wins, fund: number, mrtp: number): void { this.Slotx.Spawn(wins, fund, mrtp); }
    Prepare(): void { this.Slotx.Prepare(); }
    Apply(wins: Wins): void { this.Slotx.Apply(wins); }
    GetGain(): number { return this.Slotx.GetGain(); }
    SetGain(gain: number): Promise<void> { return this.Slotx.SetGain(gain); }
    GetBet(): number { return this.Slotx.GetBet(); }
    SetBet(bet: number): Promise<void> { return this.Slotx.SetBet(bet); }
    GetSel(): number { return this.Slotx.GetSel(); }
    // SetSel is special for this game
    // SetMode will use Slotx default (ErrNoFeature)


    Clone(): SlotGame {
        const newGame = new Game();
        // Clone Screen5x3 properties
        newGame.Screen = this.Screen.map(row => row.slice());
        // Clone Slotx properties
        newGame.Slotx = { ...this.Slotx };
        return newGame;
    }

    async Scanner(wins: Wins): Promise<void> {
        this.ScanLined(wins);
        // No error condition in Go version, so resolve promise
        return Promise.resolve();
    }

    ScanLined(wins: Wins): void {
        for (let li = 1; li <= this.Sel; li++) {
            const line = BetLines[li - 1];
            if (!line) continue;

            let multWild: number = 1;
            let numLine: Pos = 5 as Pos; // Assume 5 symbols match initially
            const firstSymbolOnLine = this.LY(1 as Pos, line); // LY is 1-indexed for position on line

            for (let x: Pos = 2 as Pos; x <= 5; x++) {
                const currentSymbol = this.LY(x, line);
                if (currentSymbol === WILD_SYMBOL) {
                    multWild = 2;
                } else if (currentSymbol !== firstSymbolOnLine) {
                    numLine = (x - 1) as Pos;
                    break;
                }
            }

            // Symbol indices are 1-based in LinePay (e.g. Cherry is symbol 2, index 1 in LinePay)
            // numLine is 1-based for number of symbols (e.g. 2 symbols, index 1 in LinePay's sub-array)
            const payTableSymbolIndex = firstSymbolOnLine -1; // Adjust for 0-indexed array
            const payTableNumIndex = numLine -1; // Adjust for 0-indexed array

            if (payTableSymbolIndex < 0 || payTableSymbolIndex >= LinePay.length || !LinePay[payTableSymbolIndex] || payTableNumIndex <0) {
                continue;
            }
            const pay = LinePay[payTableSymbolIndex][payTableNumIndex];

            if (pay > 0) {
                let multMode: number = 1;
                let freeSpinsAwarded: number = 0;

                if (this.FSR && this.FSR > 0) { // Free Spin Mode
                    multMode = 2;
                    if (LineFreespinBon[payTableSymbolIndex] && payTableNumIndex < LineFreespinBon[payTableSymbolIndex].length) {
                         freeSpinsAwarded = LineFreespinBon[payTableSymbolIndex][payTableNumIndex] || 0;
                    }
                } else { // Regular Mode
                     if (LineFreespinReg[payTableSymbolIndex] && payTableNumIndex < LineFreespinReg[payTableSymbolIndex].length) {
                        freeSpinsAwarded = LineFreespinReg[payTableSymbolIndex][payTableNumIndex] || 0;
                    }
                }

                const winItem: WinItem = {
                    Pay: this.Bet * pay,
                    Mult: multWild * multMode,
                    Sym: firstSymbolOnLine,
                    Num: numLine,
                    Line: li,
                    XY: line.slice(0, numLine) as Linex, // Get only the winning part of the line
                    Free: freeSpinsAwarded,
                };
                // Ensure XY has CopyL if needed by consumers, though fruitshop itself doesn't use it post-creation.
                // For safety, we can wrap it:
                const xyProper = NewLinex(numLine);
                for(let i=0; i<numLine; i++) { xyProper[i] = line[i]; }
                winItem.XY = xyProper;

                wins.push(winItem);
            }
        }
    }

    Spin(mrtp: number): void {
        if (Object.keys(ReelsMap).length === 0) {
            console.error("ReelsMap is empty. Cannot perform spin. Please load reels first.");
            throw new Error("ReelsMap is empty, cannot select reels for spin.");
        }
        const { reels } = FindClosest(ReelsMap, mrtp);
        super.ReelSpin(reels); // Calls Screen5x3.ReelSpin
    }

    SetSel(sel: number): Promise<void> {
        // Fruitshop does not allow changing selected lines, as per original `slot.ErrNoFeature`
        return Promise.reject(ErrNoFeature);
    }

    SetMode(n: number): Promise<void> {
        // Fruitshop does not have special modes to set via SetMode
        return Promise.reject(ErrNoFeature);
    }
}

// Helper function to create a new game instance, similar to NewGame in Go.
export function NewGame(): Game {
    return new Game();
}

// Function to populate ReelsMap, e.g. from a JSON/YAML loaded configuration
// This is a stub and would be replaced by actual loading logic.
export function LoadReels(reelsData: any) {
    const parsedReelsMap: { [key: string]: Reels5x } = {};
    for (const rtpKey in reelsData) {
        // Assuming reelsData[rtpKey] is an array of 5 reel strips (arrays of Sym)
        // e.g., "96.71": [[1,2,3,...], [4,5,6,...], ..., [7,8,9,...]]
        const reelStrips = reelsData[rtpKey] as [Sym[], Sym[], Sym[], Sym[], Sym[]];
        if (Array.isArray(reelStrips) && reelStrips.length === 5 && reelStrips.every(Array.isArray)) {
             parsedReelsMap[rtpKey] = new Reels5x(reelStrips);
        } else {
            console.warn(`Invalid reel data for RTP key ${rtpKey}:`, reelStrips);
        }
    }
    ReelsMap = parsedReelsMap;
}
