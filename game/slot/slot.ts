// Basic types from slot.go
export type Sym = number; // byte in Go, using number in TS
export type Pos = number; // int8 in Go, using number in TS

// Placeholder for Linex, used in WinItem
export interface Linex extends Array<[Pos, Pos]> {
    CopyL: (num: Pos) => Linex;
}

export const NewLinex = (length: number = 0): Linex => {
    const arr = new Array(length) as Linex;
    arr.CopyL = (num: Pos): Linex => {
        const newArr = arr.slice(0, num) as Linex;
        newArr.CopyL = arr.CopyL; // maintain the method
        return newArr;
    }
    return arr;
}


// WinItem describes win on each line or by scatters.
export interface WinItem {
  Pay?: number;
  Mult?: number;
  Sym?: Sym;
  Num?: Pos;
  Line?: number;
  XY: Linex;
  Free?: number;
  BID?: number;
  Bon?: any;
  JID?: number;
  Jack?: number;
}

// Wins is full list of wins by all lines and scatters for some spin.
export interface Wins extends Array<WinItem> {
  Gain: () => number;
  Jackpot: () => number;
  Reset: () => void;
}

export const NewWins = (): Wins => {
    const winsArray = [] as WinItem[] as Wins;

    winsArray.Gain = function(): number {
        let sum = 0;
        for (const wi of this) {
            sum += (wi.Pay || 0) * (wi.Mult || 1);
        }
        return sum;
    };

    winsArray.Jackpot = function(): number {
        let sum = 0;
        for (const wi of this) {
            sum += wi.Jack || 0;
        }
        return sum;
    };

    winsArray.Reset = function(): void {
        this.length = 0;
    };

    return winsArray;
}

// Reels interface
export interface Reels {
  Cols: () => number;
  Reel: (col: Pos) => Sym[];
  Reshuffles: () => number;
  toString: () => string;
}

// Reels for 5-reels slots.
export class Reels5x implements Reels {
  reels: [Sym[], Sym[], Sym[], Sym[], Sym[]];

  constructor(reelsData: [Sym[], Sym[], Sym[], Sym[], Sym[]]) {
    this.reels = reelsData;
  }

  Cols(): number {
    return 5;
  }

  Reel(col: Pos): Sym[] {
    if (col < 1 || col > 5) {
      throw new Error("Column out of bounds");
    }
    return this.reels[col - 1];
  }

  Reshuffles(): number {
    return (
      this.reels[0].length *
      this.reels[1].length *
      this.reels[2].length *
      this.reels[3].length *
      this.reels[4].length
    );
  }

  toString(): string {
    return `[${this.reels[0].length}, ${this.reels[1].length}, ${this.reels[2].length}, ${this.reels[3].length}, ${this.reels[4].length}]`;
  }
}

// Errors (simplified)
export const ErrNoFeature = new Error("feature not available");
export const ErrBetEmpty = new Error("bet is empty");
export const ErrDisabled = new Error("feature is disabled");
export const ErrNoLineset = new Error("lines set is empty");
export const ErrLinesetOut = new Error("lines set is out of range bet lines");


// Slotx is base struct for all slot games
export class Slotx {
  Sel: number; // selected bet lines
  Bet: number; // bet value
  Gain?: number; // gain for double up games
  FSN?: number; // free spin number
  FSR?: number; // free spin remains

  constructor(sel: number, bet: number) {
    this.Sel = sel;
    this.Bet = bet;
    this.Gain = 0;
    this.FSN = 0;
    this.FSR = 0;
  }

  Cost(): { cost: number, isJp: boolean } {
    return { cost: this.Bet * this.Sel, isJp: false };
  }

  Free(): boolean {
    return (this.FSR || 0) !== 0;
  }

  Spawn(wins: Wins, fund: number, mrtp: number): void {
    // Default implementation, can be overridden
  }

  Prepare(): void {
    // Default implementation, can be overridden
  }

  Apply(wins: Wins): void {
    if ((this.FSR || 0) !== 0) {
      this.Gain = (this.Gain || 0) + wins.Gain();
      this.FSN = (this.FSN || 0) + 1;
    } else {
      this.Gain = wins.Gain();
      this.FSN = 0;
    }

    if ((this.FSR || 0) > 0) {
      this.FSR = (this.FSR || 0) - 1;
    }
    for (const wi of wins) {
      if ((wi.Free || 0) > 0) {
        this.FSR = (this.FSR || 0) + (wi.Free || 0);
      }
    }
  }

  GetGain(): number {
    return this.Gain || 0;
  }

  SetGain(gain: number): Promise<void> {
    this.Gain = gain;
    return Promise.resolve();
  }

  GetBet(): number {
    return this.Bet;
  }

  SetBet(bet: number): Promise<void> {
    if (bet <= 0) {
      return Promise.reject(ErrBetEmpty);
    }
    if (bet === this.Bet) {
      return Promise.resolve();
    }
    if ((this.FSR || 0) !== 0) {
      return Promise.reject(ErrDisabled);
    }
    this.Bet = bet;
    return Promise.resolve();
  }

  GetSel(): number {
    return this.Sel;
  }

  SetSel(sel: number, bln?: number): Promise<void> {
    // bln is max lines, specific games will implement this fully
    if (sel < 1) {
        return Promise.reject(ErrNoLineset);
    }
    if (bln && sel > bln) {
        return Promise.reject(ErrLinesetOut);
    }
    if (sel === this.Sel) {
        return Promise.resolve();
    }
    if ((this.FSR || 0) !== 0) {
        return Promise.reject(ErrDisabled);
    }
    this.Sel = sel;
    return Promise.resolve();
  }

  SetMode(n: number): Promise<void> {
    return Promise.reject(ErrNoFeature);
  }
}

// Screen interface (placeholder, as it's embedded in Game)
// In Go, this is likely an array of arrays representing the game grid.
// type Screen [][]Sym
export interface Screen {
  // Example: GetSymbol(row: Pos, col: Pos): Sym;
  // Example: SetSymbol(row: Pos, col: Pos, sym: Sym): void;
  // For Screen5x3, it would be a 3x5 grid.
  Screen: Sym[][]; // 3 rows, 5 columns
  LY(y: Pos, line: Line): Sym; // Get symbol by Y position on a line
  ReelSpin(reels: Reels5x): void; // Spin the reels
}

// Screen5x3 embedding (simplified representation)
// The Go code uses `yaml:",inline"` which means its fields are part of the parent struct.
// We'll represent this by having the Game class implement Screen methods or have a Screen property.
export class Screen5x3 implements Screen {
    Screen: Sym[][]; // 3 rows, 5 columns

    constructor() {
        this.Screen = Array(3).fill(null).map(() => Array(5).fill(0 as Sym));
    }

    // Get symbol by Y position (row) on a specific bet line configuration.
    // line: Line object that maps a line position to screen X,Y coordinates.
    // y: The position on the line (1-indexed).
    LY(y: Pos, line: Line): Sym {
        // Assuming Line is an array of [row, col] for each position on the line
        // And y is 1-indexed for the position on that line path
        if (y < 1 || y > line.length) {
            throw new Error("Invalid position on line");
        }
        const [row, col] = line[y-1]; // line[0] is the first symbol of the line.
        if (row < 0 || row >= 3 || col < 0 || col >= 5) {
            throw new Error("Line coordinate out of screen bounds");
        }
        return this.Screen[row][col];
    }

    ReelSpin(reels: Reels5x): void {
        for (let col: Pos = 0; col < 5; col++) {
            const reelStrip = reels.Reel(col + 1 as Pos);
            const stopPos = Math.floor(Math.random() * reelStrip.length);
            for (let row: Pos = 0; row < 3; row++) {
                this.Screen[row][col] = reelStrip[(stopPos + row) % reelStrip.length];
            }
        }
    }
}


// SlotGame interface
export interface SlotGame extends Slotx, Screen {
  Clone(): SlotGame;
  Scanner(wins: Wins): Promise<void>; // Errors handled by Promise rejection
  // Cost, Free, Spin, Spawn, Prepare, Apply, GetGain, SetGain, GetBet, SetBet, GetSel, SetSel, SetMode are in Slotx
  // LY, ReelSpin are in Screen
}

// BetLine definition based on BetLinesNetEnt5x3 usage (array of [row, col] tuples)
// Each inner array represents [row, col] for a position on the line.
// e.g., [[0,0], [0,1], [0,2], [0,3], [0,4]] would be the top row.
export type Line = Array<[Pos, Pos]>; // Array of [row, col] coordinates
export const BetLinesNetEnt5x3: Line[] = [
    // Line 1 (middle row)
    [[1,0],[1,1],[1,2],[1,3],[1,4]],
    // Line 2 (top row)
    [[0,0],[0,1],[0,2],[0,3],[0,4]],
    // Line 3 (bottom row)
    [[2,0],[2,1],[2,2],[2,3],[2,4]],
    // Line 4 (diagonal top-left to bottom-right)
    [[0,0],[1,1],[2,2],[1,3],[0,4]],
    // Line 5 (diagonal bottom-left to top-right)
    [[2,0],[1,1],[0,2],[1,3],[2,4]],
    // Line 6
    [[0,0],[0,1],[1,2],[2,3],[2,4]],
    // Line 7
    [[2,0],[2,1],[1,2],[0,3],[0,4]],
    // Line 8
    [[1,0],[0,1],[0,2],[0,3],[1,4]],
    // Line 9
    [[1,0],[2,1],[2,2],[2,3],[1,4]],
    // Line 10
    [[0,0],[1,1],[1,2],[1,3],[0,4]],
    // Line 11
    [[2,0],[1,1],[1,2],[1,3],[2,4]],
    // Line 12
    [[1,0],[0,1],[1,2],[2,3],[1,4]],
    // Line 13
    [[1,0],[2,1],[1,2],[0,3],[1,4]],
    // Line 14
    [[0,0],[1,1],[2,2],[2,3],[2,4]],
    // Line 15
    [[2,0],[1,1],[0,2],[0,3],[0,4]],
];


// Placeholder for slot.Stat
// This seems to be a statistics accumulator.
export interface Stat {
  Count: () => number;
  LineRTP: (cost: number) => number;
  ScatRTP: (cost: number) => number;
  FreeCount: () => number;
  FreeHits: () => number;
  // Add other methods as they are discovered or needed
  AddWin: (winItem: WinItem, cost: number) => void; // Example method
  AddSpin: (isFreeSpin: boolean, hasFreeHit: boolean) => void; // Example method
}

export const NewStat = (): Stat => {
    let spinCount = 0;
    let totalLinePay = 0;
    let totalScatPay = 0;
    let freeSpinAwardedCount = 0;
    let freeSpinHitCount = 0;

    return {
        Count: () => spinCount,
        LineRTP: (cost: number) => (cost > 0 && spinCount > 0 ? (totalLinePay / (spinCount * cost)) * 100 : 0),
        ScatRTP: (cost: number) => (cost > 0 && spinCount > 0 ? (totalScatPay / (spinCount * cost)) * 100 : 0),
        FreeCount: () => freeSpinAwardedCount,
        FreeHits: () => freeSpinHitCount,
        AddWin: (winItem: WinItem, cost: number) => {
            const gain = (winItem.Pay || 0) * (winItem.Mult || 1);
            if (winItem.Line && winItem.Line > 0) { // Assuming line wins
                totalLinePay += gain;
            } else { // Assuming scatter wins or other non-line wins
                totalScatPay += gain;
            }
            if (winItem.Free && winItem.Free > 0) {
                freeSpinAwardedCount += winItem.Free;
            }
        },
        AddSpin: (isFreeSpin: boolean, hasFreeHit: boolean) => {
            if (!isFreeSpin) { // Only count paid spins for RTP base
                spinCount++;
            }
            if (hasFreeHit) {
                freeSpinHitCount++;
            }
        }
    };
};

// Placeholder for slot.FindClosest
// Finds the Reels5x configuration closest to the target RTP.
export const FindClosest = (reelsMap: { [rtp: string]: Reels5x }, mrtp: number): { reels: Reels5x, rtp: number } => {
  let closestRtp = -1000;
  let selectedReels: Reels5x | undefined = undefined;

  for (const rtpStr in reelsMap) {
    const rtp = parseFloat(rtpStr);
    if (Math.abs(mrtp - rtp) < Math.abs(mrtp - closestRtp)) {
      closestRtp = rtp;
      selectedReels = reelsMap[rtpStr];
    }
  }
  if (!selectedReels) {
      // Fallback or error if no reels are found, though the Go code implies one will always be found.
      // Picking the first available one if logic above somehow fails (should not happen with mrtp logic)
      const firstKey = Object.keys(reelsMap)[0];
      if (!firstKey) throw new Error("ReelsMap is empty");
      selectedReels = reelsMap[firstKey];
      closestRtp = parseFloat(firstKey);
  }
  return { reels: selectedReels!, rtp: closestRtp };
};

// Placeholder for slot.ScanReels5x
// Scans reels for wins. This is a complex function.
// The actual implementation would iterate through all possible reel stops.
export const ScanReels5x = (
  ctx: any, // context, assuming 'any'
  s: Stat,
  g: SlotGame,
  reels: Reels5x,
  calc: (writer: { write: (msg: string) => void }) => number
): number => {
  // This is a simplified simulation. A full scan is computationally intensive.
  // It would involve iterating g.Reshuffles() times.
  // For now, we'll simulate a number of spins and update stats.
  const numSimulations = 100000; // Arbitrary number of simulations
  const costInfo = g.Cost();
  const baseCost = costInfo.cost;

  for (let i = 0; i < numSimulations; i++) {
    g.ReelSpin(reels); // Spin the reels on a copy of the game
    const wins = NewWins();
    g.Scanner(wins); // Calculate wins (asynchronously, but we'll await it conceptually)

    let hasFreeHitOnSpin = false;
    for (const win of wins) {
        s.AddWin(win, baseCost);
        if (win.Free && win.Free > 0) {
            hasFreeHitOnSpin = true;
        }
    }
    s.AddSpin(g.Free(), hasFreeHitOnSpin); // Record if it was a free spin or awarded free spins
  }

  // The actual RTP calculation would be derived from 's' after all reshuffles.
  // The 'calc' function is used to format and print stats.
  let calculatedRtp = 0;
  const writer = { write: (msg: string) => console.log(msg) }; // Basic writer
  if (s.Count() > 0) { // Ensure spins were counted to avoid division by zero
    calculatedRtp = calc(writer);
  } else {
    console.warn("ScanReels5x: Zero spins counted in Stat object, cannot calculate RTP via calc callback.");
  }

  return calculatedRtp; // Return the RTP calculated by the 'calc' function
};

// Placeholder for slot.ReadMap and slot.ReadObj
// These would typically involve parsing YAML or JSON.
// For now, we'll assume they return a pre-loaded object or a typed map.
export const ReadMap = <T>(data: any /* normally []byte for YAML */): { [key: string]: T } => {
  // In a real scenario, this would parse the YAML data from []byte.
  // For FruitShop, fruitshop_reel.yaml contains RTP-keyed reel configurations.
  // We'll assume 'data' is already in the correct map format for TS.
  // e.g. { "96.0": new Reels5x(...), "97.0": new Reels5x(...) }
  if (typeof data === 'object' && data !== null) {
    return data as { [key: string]: T };
  }
  console.warn("ReadMap received non-object data, returning empty map. Data:", data);
  return {};
};

export const ReadObj = <T>(data: any /* normally []byte for YAML */): T => {
  // Similar to ReadMap, but for a single object.
  if (typeof data === 'object' && data !== null) {
    return data as T;
  }
  // This is problematic as we don't know what T is to return a default.
  // Throwing an error or returning undefined might be better.
  console.warn("ReadObj received non-object data. Data:", data);
  return undefined as any as T; // Or throw error
};

// Context placeholder (if needed, though ScanReels5x simplifies it)
export interface Context {
  Err: () => Error | null;
  // Other context methods if used by the game logic
}
