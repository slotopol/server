import {
  Sym,
  Pos,
  Linex,
  SlotGame,
  Slotx,
  Screen5x3,
  Wins,
  WinItem,
  Reels5x,
  findClosestReels,
  readObj,
  ErrNoLineset,
  ErrLinesetOut,
  ErrDisabled,
} from '../../slot/types';
import { ReelsBon, ReelsMap } from './flowers_reels';

// Constants from flowers_rule.go
// Symbol IDs (1-indexed as in Go)
export const wild: Sym = 1;
export const scat: Sym = 2;
export const scat2: Sym = 3; // scatter2 in Go

// Lined payment.
export const LinePay: number[][] = [
  /* 0: wild */ [0, 0, 250, 1000, 5000],
  /* 1: scatter (empty in Go, used for scat2) */ [],
  /* 2: scatter2 (empty in Go) */ [],
  /* 3: red */ [0, 0, 20, 40, 160, 250, 400, 600, 1000, 2000],
  /* 4: red2 */ [],
  /* 5: yellow */ [0, 0, 15, 35, 140, 225, 350, 550, 900, 1800],
  /* 6: yellow2 */ [],
  /* 7: green */ [0, 0, 15, 30, 120, 200, 300, 500, 800, 1600],
  /* 8: green2 */ [],
  /* 9: pink */ [0, 0, 10, 25, 100, 175, 250, 450, 700, 1400],
  /* 10: pink2 */ [],
  /* 11: blue */ [0, 0, 10, 20, 80, 150, 200, 400, 600, 1200],
  /* 12: blue2 */ [],
  /* 13: ace */ [0, 0, 5, 20, 200],
  /* 14: king */ [0, 0, 5, 20, 150],
  /* 15: queen */ [0, 0, 5, 15, 125],
  /* 16: jack */ [0, 0, 5, 15, 100],
];
// Note: Go's LinePay is [17][10]float64. TS uses 0-indexed arrays.
// Access in Go: LinePay[wild-1][numw-1]
// Access in TS: LinePay[wildSymbolIdAsArrayIndex][numberOfSymbolsAsArrayIndex]
// We'll adjust symbol IDs when accessing these arrays. For example, wild (1) becomes index 0.

export const DoubleSym: Sym[] = [
  // 0-indexed: symbol ID (1-17) maps to index (0-16)
  /* 1 wild   */ 0,
  /* 2 scat   */ 0,
  /* 3 scat2  */ 2, // scat2 (ID 3) points to original scatter (ID 2)
  /* 4 red    */ 0,
  /* 5 red2   */ 4, // red2 (ID 5) points to original red (ID 4)
  /* 6 yellow */ 0,
  /* 7 yellow2*/ 6,
  /* 8 green  */ 0,
  /* 9 green2 */ 8,
  /* 10 pink  */ 0,
  /* 11 pink2 */ 10,
  /* 12 blue  */ 0,
  /* 13 blue2 */ 12,
  /* 14 ace   */ 0,
  /* 15 king  */ 0,
  /* 16 queen */ 0,
  /* 17 jack  */ 0,
];

// Scatters payment.
export const ScatPay: number[] = [0, 0, 2, 2, 2, 2, 4, 10]; // Max 8 scatters (index 0-7 for 1-8 scatters)
// Go: ScatPay[count-1]

// Scatter freespins table
export const ScatFreespin: number[] = [0, 0, 0, 10, 15, 20, 25, 30]; // Max 8 scatters
// Go: ScatFreespin[count-1]

// Bet lines (1-indexed in comments, but 0-indexed in use for reel positions)
// Values are Y-positions (1, 2, or 3)
export const BetLines: Linex[] = [
  [2, 2, 2, 2, 2], // 1
  [1, 1, 1, 1, 1], // 2
  [3, 3, 3, 3, 3], // 3
  [1, 2, 3, 2, 1], // 4
  [3, 2, 1, 2, 3], // 5
  [1, 1, 2, 1, 1], // 6
  [3, 3, 2, 3, 3], // 7
  [2, 3, 3, 3, 2], // 8
  [2, 1, 1, 1, 2], // 9
  [2, 1, 2, 1, 2], // 10
  [2, 3, 2, 3, 2], // 11
  [1, 2, 1, 2, 1], // 12
  [3, 2, 3, 2, 3], // 13
  [2, 2, 1, 2, 2], // 14
  [2, 2, 3, 2, 2], // 15
  [1, 2, 2, 2, 1], // 16
  [3, 2, 2, 2, 3], // 17
  [1, 3, 1, 3, 1], // 18
  [3, 1, 3, 1, 3], // 19
  [1, 3, 3, 3, 1], // 20
  [3, 1, 1, 1, 3], // 21
  [1, 1, 3, 1, 1], // 22
  [3, 3, 1, 3, 3], // 23
  [1, 3, 2, 1, 3], // 24
  [3, 1, 2, 3, 1], // 25
  [2, 1, 3, 1, 2], // 26
  [2, 3, 1, 3, 2], // 27
  [1, 2, 3, 3, 3], // 28
  [3, 2, 1, 1, 1], // 29
  [2, 1, 2, 3, 2], // 30
];

export class FlowersGame implements SlotGame, Slotx {
  public screen: Screen5x3;
  // Slotx properties
  public sel: number;
  public bet: number;
  public gain: number;
  public fsn: number; // free spin number
  public fsr: number; // free spin remains

  constructor() {
    this.screen = new Screen5x3();
    this.sel = BetLines.length;
    this.bet = 1;
    this.gain = 0;
    this.fsn = 0;
    this.fsr = 0;
  }

  clone(): FlowersGame {
    const newGame = new FlowersGame();
    newGame.screen = new Screen5x3(); // Ensure screen is a new instance
    // Deep copy screen state if necessary, for now new screen is fine
    for (let i = 0; i < 5; i++) {
        for (let j = 0; j < 3; j++) {
            newGame.screen.scr[i][j] = this.screen.scr[i][j];
        }
    }
    newGame.sel = this.sel;
    newGame.bet = this.bet;
    newGame.gain = this.gain;
    newGame.fsn = this.fsn;
    newGame.fsr = this.fsr;
    return newGame;
  }

  cost(): { cost: number; hasJackpotRate: boolean } {
    return { cost: this.bet * this.sel, hasJackpotRate: false };
  }

  free(): boolean {
    return this.fsr > 0;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  spawn(wins: Wins, fund: number, rtp: number): void {
    // No specific spawn logic in flowers_rule.go Slotx implementation
  }

  prepare(): void {
    // No specific prepare logic in flowers_rule.go Slotx implementation
  }

  apply(wins: Wins): void {
    if (this.fsr > 0) {
      this.gain += wins.reduce((acc, w) => acc + (w.pay || 0) * (w.mult || 1), 0);
      this.fsn++;
    } else {
      this.gain = wins.reduce((acc, w) => acc + (w.pay || 0) * (w.mult || 1), 0);
      this.fsn = 0;
    }

    if (this.fsr > 0) {
      this.fsr--;
    }
    for (const wi of wins) {
      if (wi.free && wi.free > 0) {
        this.fsr += wi.free;
      }
    }
  }

  getGain(): number {
    return this.gain;
  }

  setGain(gain: number): Error | null {
    this.gain = gain;
    return null;
  }

  getBet(): number {
    return this.bet;
  }

  setBet(bet: number): Error | null {
    if (bet <= 0) return new Error("Bet must be positive."); // Equivalent to ErrBetEmpty
    if (this.fsr > 0) return new Error("Cannot change bet during free spins."); // Equivalent to ErrDisabled
    this.bet = bet;
    return null;
  }

  getSel(): number {
    return this.sel;
  }

  setSel(sel: number): Error | null {
    if (sel < 1) return ErrNoLineset;
    if (sel > BetLines.length) return ErrLinesetOut;
    if (this.fsr > 0) return ErrDisabled;
    this.sel = sel;
    return null;
  }

  // Slotx's SetSelNum, aliased to avoid conflict with SlotGame's SetSel
  setSelNum(sel: number, bln: number): Error | null {
    if (sel < 1) return ErrNoLineset;
    if (sel > bln) return ErrLinesetOut;
    if (this.fsr > 0) return ErrDisabled;
    this.sel = sel;
    return null;
  }


  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  setMode(mode: number): Error | null {
    return new Error("Feature not available."); // Equivalent to ErrNoFeature
  }

  scanner(wins: Wins): Error | null {
    wins.length = 0; // Clear previous wins, similar to Go's wins.Reset() if it were a pointer
    this.scanLined(wins);
    this.scanScatters(wins);
    return null;
  }

  scanLined(wins: Wins): void {
    for (let li = 0; li < this.sel; li++) { // Iterate up to selected lines
      const lineDef = BetLines[li]; // 0-indexed line definition

      let numw: Pos = 0;
      let numl: Pos = 0;
      let syml: Sym = 0;
      let x: Pos; // Represents reel number 1-5

      for (x = 1; x <= 5; x++) {
        const sx = this.screen.ly(x, lineDef); // ly expects 1-indexed x and lineDef

        if (sx === wild) {
          if (syml === 0) { // Only count wilds if no other symbol started the line
            numw++;
          }
          // If a symbol line already started, wilds can substitute for that symbol
          // The Go logic is a bit complex here:
          // It seems to prioritize non-wilds first. If wild starts, numw counts.
          // If non-wild starts, syml is set. Wilds encountered later seem to increment numl if they match syml.
          // This needs careful translation. The Go code:
          // if sx == wild { if syml == 0 { numw++ } } else if sd := DoubleSym[sx-1]; syml == 0 { ... }
          // This implies wilds are only counted for `numw` if they form their own line from the left.
          // If `syml` is already set, a wild `sx` would typically extend `numl` if `wild` can substitute `syml`.
          // The current Go code seems to break if wild is encountered after syml is set and sx != syml.
          // Let's stick to the Go logic: if `sx` is wild, it only contributes to `numw` if `syml` is 0.
          // Otherwise, if `syml` is set, a wild does not break the line but also doesn't extend `numl` directly in this loop for `syml`.
          // The win selection logic below handles choosing between wild pay and symbol pay.
        } else { // Not a wild
            const doubleSymEquiv = DoubleSym[sx -1]; // sx is 1-indexed symbol ID
            if (syml === 0) { // First non-wild symbol on the line
                if (doubleSymEquiv > 0) { // It's a double symbol
                    syml = doubleSymEquiv; // Use the single equivalent for payline
                    numl++; // Counts as two
                } else {
                    syml = sx; // It's a single symbol
                }
            } else if (doubleSymEquiv > 0 && doubleSymEquiv === syml) { // Double symbol matching current line symbol
                numl++; // Counts as two
            } else if (sx === syml) { // Single symbol matching current line symbol
                // numl is already incremented by one below
            } else { // Symbol mismatch
                break; // End of current symbol line
            }
        }
        // Common logic for non-wilds and wilds that extend a non-wild line (if that were the case)
        if (sx !== wild) { // Only increment numl for actual symbols or their double versions
            numl++;
        } else if (syml !== 0) { // Wild appearing after a symbol line has started
            // If wild substitutes for syml, increment numl.
            // The original Go code doesn't explicitly do this here, it relies on payw vs payl.
            // For now, let's assume wild extends the line of `syml`.
            numl++;
        }


      } // End of reel scan for one line

      const currentReelScanLength = x-1; // How many reels were successfully matched for syml or wild

      let payw = 0;
      // Go: LinePay[wild-1][numw-1]
      if (numw >= 3) { // Wilds pay from 3
        payw = LinePay[wild - 1][numw - 1] || 0;
      }

      let payl = 0;
      // Go: LinePay[syml-1][numl-1]
      if (numl >= 3 && syml > 0) { // Other symbols pay from 3
        // syml is 1-indexed symbol ID
        const symIndex = syml -1;
        if (LinePay[symIndex] && numl-1 < LinePay[symIndex].length) {
             payl = LinePay[symIndex][numl - 1] || 0;
        }
      }

      let winAdded = false;
      if (payl > payw) {
        const winItem: WinItem = {
          pay: this.bet * payl,
          mult: this.fsr > 0 ? 3 : 1,
          sym: syml,
          num: numl,
          line: li + 1, // 1-indexed line number for output
          xy: lineDef.slice(0, currentReelScanLength) as Linex, // Show only matched part of line
        };
        wins.push(winItem);
        winAdded = true;
      } else if (payw > 0) {
        const winItem: WinItem = {
          pay: this.bet * payw,
          mult: this.fsr > 0 ? 3 : 1,
          sym: wild,
          num: numw,
          line: li + 1, // 1-indexed line number for output
          xy: lineDef.slice(0, numw) as Linex, // Wild line length
        };
        wins.push(winItem);
        winAdded = true;
      }
    }
  }

  scatNumDbl(): Pos {
    let n: Pos = 0;
    for (let x = 1; x <= 5; x++) { // Iterate through reels (1 to 5)
      let foundOnReel = false;
      let isDouble = false;
      for (let y = 1; y <= 3; y++) { // Iterate through rows (1 to 3)
        const sym = this.screen.at(x,y);
        if (sym === scat) {
          foundOnReel = true;
          break;
        }
        if (sym === scat2) {
          foundOnReel = true;
          isDouble = true;
          break;
        }
      }
      if (foundOnReel) {
        n += isDouble ? 2 : 1;
      }
    }
    return n;
  }

  scanScatters(wins: Wins): void {
    const count = this.scatNumDbl();
    if (count >= 3) { // Scatters pay from 3
      // Go: ScatPay[count-1], ScatFreespin[count-1]
      // Ensure count-1 is a valid index
      const payIndex = Math.min(count - 1, ScatPay.length - 1);
      const fsIndex = Math.min(count - 1, ScatFreespin.length - 1);

      const pay = ScatPay[payIndex] || 0;
      const fs = ScatFreespin[fsIndex] || 0;

      const winItem: WinItem = {
        pay: this.bet * this.sel * pay, // Scatter pay is often multiplied by total bet (bet * lines)
        mult: this.fsr > 0 ? 3 : 1,
        sym: scat, // Original scatter symbol
        num: count,
        xy: this.screen.scatPos(scat).concat(this.screen.scatPos(scat2)), // Combine positions of both scat types
        free: fs,
      };
      wins.push(winItem);
    }
  }

  spin(mrtp: number): void {
    let currentReels: Reels5x;
    if (this.fsr === 0) {
      const reels = findClosestReels(ReelsMap, mrtp);
      if (!reels) {
        console.error("No reels found for RTP:", mrtp, "Using bonus reels as fallback.");
        currentReels = readObj(ReelsBon); // Fallback, though not ideal
      } else {
        currentReels = reels;
      }
    } else {
      currentReels = readObj(ReelsBon);
    }

    // Adapt Reels5x (Sym[][]) to Reels interface for reelSpin
    const reelsAdapter: Reels = {
        cols: () => 5,
        reel: (col: Pos) => currentReels[col-1], // col is 1-indexed
        reshuffles: () => {
            if (!currentReels || currentReels.length === 0) return 0;
            return currentReels.reduce((acc, r) => acc * (r?.length || 1), 1);
        },
        toString: () => currentReels.map(r => r?.length || 0).join(', ')
    };
    this.screen.reelSpin(reelsAdapter);
  }
}

export function NewFlowersGame(): FlowersGame {
  return new FlowersGame();
}
