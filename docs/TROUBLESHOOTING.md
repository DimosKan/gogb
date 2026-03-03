# Troubleshooting Guide

## Common Problems and Solutions

This guide maps symptoms to likely causes. When your emulator misbehaves, find your symptom below.

---

## Screen Issues

### Completely Black Screen (but should show something)

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| LCD not enabled | Check if bit 7 of 0xFF40 is set | Write 0x91 to 0xFF40 (LCDC) |
| Wrong palette | Check BGP (0xFF47) | Use 0xE4 for standard, 0x03 for inverted |
| VRAM not written | Inspect VRAM at 0x8000 | Verify tile data was copied |
| GPU not called | Add debug print in Render() | Call GPU.Render() each frame |

**Step to Review:** Step 1

### Completely White Screen

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Palette is 0x00 | Check BGP register | Color 0 = white means blank is white |
| LCD disabled | Bit 7 of LCDC is 0 | Enable LCD with 0x91 |

**Step to Review:** Step 1

### Checkerboard Pattern Everywhere

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Writing to tile 0 only | Check tile map at 0x9800 | Tile map defaults to 0, so you only see tile 0 |
| Tile addressing wrong | Verify tile index calculation | Each tile is 16 bytes, tileAddr = tileIndex * 16 |

**Step to Review:** Step 2

### Garbage/Random Tiles on Screen

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| VRAM not cleared | Real hardware has garbage in VRAM | Clear tile map (0x9800-0x9BFF) before drawing |
| Tile map not set | Check what values are at 0x9800 | Write correct tile indices to tile map |
| Wrong tile data address | Verify LCDC bit 4 | Bit 4 = 0 uses 0x8800 base, bit 4 = 1 uses 0x8000 |

**Step to Review:** Step 3

### Sprites Not Visible

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Sprites not enabled | Check LCDC bit 1 | Set bit 1 of 0xFF40 |
| OBP0/OBP1 is 0x00 | Check 0xFF48/0xFF49 | Write valid palette (e.g., 0xE4) |
| Sprite off-screen | Check OAM Y position | Y < 16 or Y > 160 means off-screen |
| X position wrong | Check OAM X position | X < 8 means partially off-screen left |
| OAM not written | Inspect 0xFE00-0xFE9F | Verify sprite data was written |

**Step to Review:** Step 6

### Sprites Offset by 16 Pixels Vertically

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Forgot Y offset | OAM Y = screen Y + 16 | Subtract 16 when converting to screen coords |

**Step to Review:** Step 6

### Sprites Offset by 8 Pixels Horizontally

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Forgot X offset | OAM X = screen X + 8 | Subtract 8 when converting to screen coords |

**Step to Review:** Step 6

### Screen Flickers or Tears

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Writing VRAM during draw | Check when VRAM writes happen | Only write during VBlank (LY >= 144) |
| No VBlank waiting | Check for HALT or LY polling | Implement proper VBlank timing |

**Step to Review:** Step 7

---

## Execution Issues

### Game Crashes to 0x0000

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Little-endian reversed | Check JP a16 implementation | Low byte first: `addr = hi<<8 OR lo` |
| Unknown opcode panic | Check if unknown opcode handling | Make sure to log the offending opcode |

**Step to Review:** Step 1

### Infinite Loop (Never Halts)

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Zero flag not set by DEC | Check DEC B implementation | `gb.CPU.SetZ(gb.CPU.B == 0)` |
| JR NZ offset wrong | Check signed offset handling | Offset must be `int8`, not `uint8` |
| JR NZ condition inverted | Check condition logic | Jump when Z is NOT set |

**Step to Review:** Step 3

### PC Jumps to Wrong Address

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Byte order wrong | Print lo and hi bytes | Remember: `[lo, hi]` in memory = `hi<<8 OR lo` |
| Forgot to advance PC | Check PC after reading operands | Some instructions need `PC += 2` |
| Signed offset as unsigned | Check JR implementation | Cast to `int8` before using |

**Step to Review:** Step 1, Step 3

### Stack Corruption

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| SP not initialized | Check SP value | Initialize SP to 0xFFFE |
| Push/pop order wrong | Stack grows downward | PUSH: decrement then write; POP: read then increment |
| Wrong byte order on stack | Check CALL/RET | Push high byte first (at higher address) |

**Step to Review:** Step 7

---

## Input Issues

### Buttons Not Working

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Active-low inverted | Check JOYP read logic | 0 = pressed, 1 = not pressed |
| Wrong selection bits | Check bits 4-5 of JOYP | Bit 5 = select D-pad, Bit 4 = select buttons |
| Key mapping wrong | Check keyboard to button mapping | Verify Ebitengine key constants |

**Step to Review:** Step 4

### Input Registers Always 0xFF

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| JOYP not connected | Check MMU.Read for 0xFF00 | Return Joypad.Read() for address 0xFF00 |
| Joypad.Read() wrong | Check return value | Should return 0xCF or similar when no buttons pressed |

**Step to Review:** Step 4

---

## Timing Issues

### Animation Too Fast

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| No frame limiting | Check cycles per frame | Should be ~70224 cycles per frame |
| Running in infinite loop | Check main loop | Use Ebitengine's 60 FPS frame limiting |

**Step to Review:** Step 7

### Animation Too Slow

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Too few cycles per frame | Check CyclesPerFrame constant | Should be 70224, not less |
| Blocking in Update() | Check for slow operations | Don't do I/O or sleep in Update() |

**Step to Review:** Step 7

---

## Flag Issues

### Carry Flag Not Set

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Checking after truncation | Check flag calculation order | Calculate flags BEFORE truncating result to 8 bits |
| Wrong condition | Check the comparison | `result > 0xFF` for ADD, `a < b` for SUB |

**Step to Review:** Step 3, Step 8

### Half-Carry Flag Not Set

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Wrong nibble check | Check half-carry calculation | `(a & 0x0F) + (b & 0x0F) > 0x0F` for ADD |
| Not implemented | Some emulators skip this | Required for BCD operations (DAA) |

**Step to Review:** Step 4

### Zero Flag Not Set

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Checking wrong value | Check what value is tested | Must check RESULT, not operand |
| Set before operation | Check order | Set Z flag AFTER the operation |

**Step to Review:** Step 3

---

## Memory Issues

### Writing to ROM Area

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| No protection | Check MMU.Write | Ignore writes to 0x0000-0x7FFF (or handle MBC) |

### VRAM Writes Ignored

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Address calculation wrong | Check VRAM offset | `gb.MMU.VRAM[addr - 0x8000]` |
| Writing during LCD draw | Check LCD mode | Should only write during VBlank or Mode 0 |

**Step to Review:** Step 2

### OAM Writes Ignored

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| OAM not implemented | Check MMU regions | Need OAM array at 0xFE00-0xFE9F |
| DMA not triggering | Check 0xFF46 write handler | DMA copies 160 bytes to OAM |

**Step to Review:** Step 6

---

## MBC Issues

### Game Freezes When Switching Areas/Levels

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Wrong ROM bank | Check which bank is selected | Verify bank register masking (5 bits for MBC1) |
| Bank 0 quirk | Game writes 0x00 to bank register | Bank 0 must map to bank 1 at 0x4000 |
| ROM size masking | Bank exceeds actual ROM | Mask bank number: `bank % romBankCount` |

**Step to Review:** Step 20

### Save Files Don't Work

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| RAM not enabled | Check if 0x0A was written to 0x0000-0x1FFF | Only enable RAM when lower nibble is 0x0A |
| No battery backup | Check cartridge type byte | Only save for types with "BATTERY" (0x03, 0x13, etc.) |
| Wrong save path | Check file extension | Use .sav extension alongside .gb file |
| RAM not persisted | Check if SaveRAM is called | Save on emulator exit and periodically |

**Step to Review:** Step 20

### Pokemon Clock is Wrong

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| RTC not latched | Check latch sequence | Must write 0x00 then 0x01 to 0x6000-0x7FFF |
| Time not tracked | Check RTC implementation | Track real elapsed time with system clock |
| Halt flag ignored | Check bit 6 of 0x0C | Stop clock updates when halt is set |

**Step to Review:** Step 20

### Game Shows Garbage After Boot

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| MBC not detected | Check cartridge type at 0x0147 | Implement MBC detection from header |
| Wrong MBC type | Verify MBC selection | Use correct MBC for cartridge type |
| ROM not loaded | Check ROM array | Load entire ROM file, not just 32KB |

**Step to Review:** Step 20

---

## Common Opcode Bugs

### NOP (0x00)

| Bug | Symptom | Fix |
|-----|---------|-----|
| Forgot to return cycles | Timing off | Return 4 |
| PC not incremented | Infinite NOP | PC increments before switch |

### JP a16 (0xC3)

| Bug | Symptom | Fix |
|-----|---------|-----|
| Big-endian instead of little | Jump to wrong address | `addr = hi<<8 OR lo`, read lo first |
| PC incremented after setting | Jump to addr+2 | Don't increment PC after setting |

### JR NZ (0x20)

| Bug | Symptom | Fix |
|-----|---------|-----|
| Unsigned offset | Can only jump forward | Cast offset to `int8` |
| Offset from wrong PC | Off by 2 | Offset is relative to PC AFTER reading instruction |
| Condition inverted | Loops forever or never | Jump when Z is NOT set (NZ = Not Zero) |

### DEC B (0x05)

| Bug | Symptom | Fix |
|-----|---------|-----|
| Zero flag not set | Loop never exits | `SetZ(B == 0)` after decrement |
| N flag not set | DAA fails | Set N flag (subtraction) |

### LD (HL+), A (0x22)

| Bug | Symptom | Fix |
|-----|---------|-----|
| HL not incremented | Overwrites same byte | `SetHL(HL() + 1)` after write |
| Increment before write | Off by one | Write first, then increment |

---

## Timing and Cycle Accuracy Issues

### Blargg's instr_timing Fails

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Wrong cycle count | Compare with opcodes.json | Verify each instruction returns correct T-cycles |
| Not counting memory accesses | Check 16-bit loads | Each memory access = 4 T-cycles |
| CB prefix not counted | Check CB instructions | CB fetch + operation cycles |
| Conditional timing wrong | Check JR/JP/CALL/RET | Taken vs not-taken have different cycles |

**Step to Review:** Step 21

### Games Work but Have Subtle Glitches

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Sprites flicker wrong | PPU mode timing | Implement precise mode 2/3 timing |
| Graphics tear mid-screen | VRAM lockout missing | Return 0xFF for VRAM reads during mode 3 |
| Sound pops/clicks | APU timing | Update APU every M-cycle, not per-instruction |
| HBlank effects broken | STAT interrupt timing | Fire STAT interrupt at exact mode change |

**Step to Review:** Step 21

### PPU Test ROMs Fail

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Mode 2 not 80 cycles | Check OAM scan time | OAM scan is exactly 80 T-cycles |
| VBlank fires wrong time | Check scanline 144 | VBlank starts at beginning of scanline 144 |
| LY updates wrong | Check LY timing | LY updates at start of each scanline |
| Mode 3 too short | Check sprite handling | Mode 3 extends with sprites on scanline |

**Step to Review:** Step 21

### Timer-Dependent Games Break

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Music plays wrong speed | Timer frequency wrong | Check TAC clock select bits |
| Events fire late | Interrupt delay | Timer interrupt on TIMA overflow |
| Timing tests fail | TIMA reload delay | TIMA reloads 4 cycles after overflow |
| DIV-based timing wrong | DIV increment rate | DIV increments every T-cycle (not M-cycle) |

**Step to Review:** Step 21

### Games With Mid-Instruction Timing Bugs

Some games observe hardware state during multi-cycle instructions:

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| Race condition crashes | Bank switch during interrupt | Game bug or VBlank timing issue |
| Interrupt fires at wrong time | Per-instruction vs per-cycle updates | Update hardware every M-cycle, not after instruction |
| Memory reads wrong value | VRAM/OAM accessed during lock | Implement PPU mode-based memory lockouts |

**Known Games with Timing Sensitivity:**
- Prehistorik Man (PPU timing)
- Road Rash (HBlank effects)
- Super Mario Land 2 (VBlank/bank switch race condition)

**Step to Review:** Step 21

---

## Audio Issues

### Audio Has Severe Latency (Sounds Delayed)

This is one of the most common audio issues in emulators. Sound effects or music play noticeably after the action that triggered them.

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Sample queue growing unbounded | Print `len(sampleQueue)` - if it keeps growing, this is the problem | Limit queue to ~100ms of samples (see fix below) |
| SetBufferSize using wrong type | Check if passing integer instead of `time.Duration` | Use `player.SetBufferSize(50 * time.Millisecond)` |
| Creating new player on every sound | Check TriggerBeep/PlaySound implementation | Use `NewPlayerFromBytes` for sound effects |
| Default audio buffer too large | No explicit SetBufferSize call | Add `player.SetBufferSize(20-50 * time.Millisecond)` |

**Root Cause Explanation:**

Streaming audio emulators generate samples in a queue that the audio system consumes. If the emulator runs even slightly faster than real-time (common), samples accumulate faster than they're consumed, causing latency to grow indefinitely.

**The Fix - Queue Limiting:**

```go
// After appending samples to the queue:
apu.sampleQueue = append(apu.sampleQueue, leftSample, rightSample)

// Limit queue size to ~100ms of audio to prevent latency buildup
// 44100 Hz * 0.1 seconds * 2 channels = 8820 samples
const maxQueueSize = 8820
if len(apu.sampleQueue) > maxQueueSize {
    // Drop oldest samples to maintain low latency
    apu.sampleQueue = apu.sampleQueue[len(apu.sampleQueue)-maxQueueSize:]
}
```

**The Fix - Correct SetBufferSize:**

```go
// WRONG - This passes nanoseconds, not milliseconds!
player.SetBufferSize(SampleRate / 20 * 4)  // ~8820 nanoseconds = 0.000009ms

// CORRECT - Use time.Duration
player.SetBufferSize(50 * time.Millisecond)  // 50ms buffer
```

**The Fix - Sound Effects (non-streaming):**

```go
func (apu *APU) TriggerBeep() {
    // Use NewPlayerFromBytes for efficient sound effect playback
    player := apu.audioContext.NewPlayerFromBytes(apu.beepData)
    // Small buffer for responsive sound effects
    player.SetBufferSize(20 * time.Millisecond)
    player.Play()
}
```

**Step to Review:** Step 9, Step 13

### No Audio at All

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| NR52 master enable off | Check bit 7 of 0xFF26 | Write 0x80 to NR52 to enable sound |
| Channel not triggered | Check NR14/NR24/NR34/NR44 bit 7 | Trigger bit must be written |
| Volume is zero | Check NR12/NR22/NR32/NR42 | Set volume envelope |
| Audio context not initialized | Check InitAudio() return value | Call InitAudio() during emulator setup |
| Panning mutes channel | Check NR51 (0xFF25) | Ensure channel is routed to speakers |

**Step to Review:** Step 9

### Audio Plays Once and Never Again

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Trigger bit not re-set | Check NR24 bit 7 handling | Sound only plays when trigger bit transitions 0→1 |
| Player not recreated | Check if reusing same player | For NewPlayerFromBytes, create new player each time |
| Length counter expired | Check if length is enabled | Length counter silences channel when it expires |

**Step to Review:** Step 9, Step 13

### Audio is Distorted or Crackly

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Sample clipping | Check sample values | Keep samples in int16 range (-32768 to 32767) |
| Buffer underrun | Audio stutters/pops | Increase buffer size or generate samples faster |
| Wrong sample rate | Check audio context | Must match: 44100 Hz typically |
| Mixing overflow | Multiple channels sum too high | Scale mixed output to prevent clipping |
| High-pass filter missing | DC offset causes distortion | Apply high-pass filter to remove DC bias |

**Step to Review:** Step 13

### Notes Cut Off Too Early

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Length counter ticking wrong | Check frame sequencer rate | Length counter ticks at 256 Hz (every other frame) |
| Envelope ends too fast | Check envelope period | Envelope ticks at 64 Hz |
| Channel disabled unexpectedly | Check DAC enable logic | DAC disable silences channel |

**Step to Review:** Step 13

### Sweep Sounds Wrong (Channel 1)

| Possible Cause | How to Check | Solution |
|----------------|--------------|----------|
| Overflow not handled | Frequency exceeds 2047 | Disable channel on overflow |
| Direction wrong | Sweep goes wrong way | Bit 3: 0=increase, 1=decrease |
| Period calculation wrong | Use correct formula | `new_freq = freq ± (freq >> shift)` |

**Step to Review:** Step 13

---

## Debugging Tips

### Add Trace Logging

```go
fmt.Printf("%04X: %02X  A:%02X F:%02X B:%02X HL:%04X\n",
    gb.CPU.PC, opcode, gb.CPU.A, gb.CPU.F, gb.CPU.B, gb.CPU.HL())
```

### Compare with Reference Emulator

Run the same ROM in BGB (Windows) or another accurate emulator and compare:
- Register values at each step
- Memory contents after operations
- Which branch is taken

### Isolate the Problem

Create a minimal ROM that demonstrates just the bug:

```go
// Test just JP a16
rom[0x100] = 0xC3  // JP
rom[0x101] = 0x50  // Low byte
rom[0x102] = 0x01  // High byte (should jump to 0x0150)
rom[0x150] = 0x76  // HALT
```

### Check Your Test First

Sometimes the test is wrong, not the implementation. Verify:
- Expected values are correct
- Little-endian is handled correctly in test setup
- PC is at the right address when test starts

---

## Quick Reference: Flag Effects

| Instruction | Z | N | H | C |
|-------------|---|---|---|---|
| INC r | * | 0 | * | - |
| DEC r | * | 1 | * | - |
| ADD A, r | * | 0 | * | * |
| SUB r | * | 1 | * | * |
| AND r | * | 0 | 1 | 0 |
| OR r | * | 0 | 0 | 0 |
| XOR r | * | 0 | 0 | 0 |
| CP r | * | 1 | * | * |
| RLCA | 0 | 0 | 0 | * |

Legend: `*` = affected, `0` = reset, `1` = set, `-` = unchanged

---

## Quick Reference: Audio Buffer Sizing

| Use Case | Buffer Size | Queue Limit | Notes |
|----------|-------------|-------------|-------|
| Sound effects (beeps) | 20ms | N/A | Use NewPlayerFromBytes, no queue |
| Streaming audio | 50ms | 8820 samples (~100ms) | Balance latency vs underrun |
| Music playback | 50-100ms | 8820-17640 samples | Can tolerate more latency |

**Key Constants:**
```go
const SampleRate = 44100          // Standard audio sample rate
const maxQueueSize = 8820         // ~100ms of stereo audio at 44100Hz
                                  // Formula: SampleRate * 0.1 * 2 channels
```

**Why These Values:**
- **20ms buffer for effects:** Responsive, but may underrun if CPU is busy
- **50ms buffer for streaming:** Good balance of latency and stability
- **100ms queue limit:** Maximum acceptable latency before dropping samples
- **Dropping samples:** Better to skip old samples than accumulate latency

---

## Mooneye Cycle-Exact Test Failures

### Tests That Require Architectural Changes

The following 4 Mooneye tests fail due to fundamental timing architecture, not implementation bugs. These require T-cycle (dot-level) precision that M-cycle granular emulators cannot achieve without rewriting.

### intr_2_mode0_timing_sprites

**What it tests:** STAT Mode 0 interrupt timing with sprites on the scanline.

**Why it fails:**
- Mode 3 length is pre-computed at scanline start
- Sprite penalties should emerge dynamically from FIFO stalls
- Interrupt dispatch timing is 1-2 T-cycles off

**What would fix it:**
- True pixel FIFO state machine
- Mode 3 ends when 160 pixels are pushed, not at pre-calculated dot
- Per-sprite fetch penalty based on background fetcher state

**Debugging approach:**
```bash
# Add logging to computeMode3End()
# Compare calculated mode3End with test expectations
# Check sprite penalty calculation for edge cases (X < 8, X > 160)
```

### lcdon_timing-GS

**What it tests:** Exact T-cycle when STAT register becomes readable after LCD enable.

**Why it fails:**
- Line 0 after LCD enable has special timing
- Our loop processes dot 0 as dot 1 (increment-first design)
- STAT mode bits update at specific T-cycles within the M-cycle

**What would fix it:**
- Process-then-increment loop structure
- Line 0 starts in Mode 0, not Mode 2
- Sub-M-cycle register sampling

**Debugging approach:**
```bash
# Log STAT reads on line 0 after LCD enable
# Record exact ScanlineCycles when read occurs
# Compare with test expectations (probes T-cycles 0-3)
```

### lcdon_write_timing-GS

**What it tests:** Exact T-cycle when register writes take effect after LCD enable.

**Why it fails:**
- Same root cause as lcdon_timing-GS
- Write timing must happen at specific T-cycle offset
- Current architecture applies writes at M-cycle boundary

**What would fix it:**
- Per-instruction write timing
- Sub-M-cycle state updates
- Precise LCD enable transition handling

### stat_lyc_onoff

**What it tests:** LYC coincidence flag behavior when LCD is toggled on/off.

**Why it fails:**
- STAT line during LCD-off state handling
- Coincidence flag update timing on LCD enable
- STAT interrupt edge detection during transitions

**What would fix it:**
- Precise STAT line state during LCD off
- Immediate LYC comparison on LCD enable
- Correct rising edge detection for STAT interrupt

**Debugging approach:**
```bash
# Log STAT coincidence flag state during LCD toggle
# Track when LYC comparison fires
# Compare STAT line high/low transitions with test expectations
```

### Summary: Why These Tests Are Different

| Test Type | Granularity Needed | Current Architecture |
|-----------|-------------------|---------------------|
| Most Mooneye tests | M-cycle | M-cycle (pass) |
| These 4 tests | T-cycle | M-cycle (fail) |

**Passing these tests requires architectural changes documented in Chapter 26.**

---

*When in doubt, write a test that isolates the behavior you're debugging.*
