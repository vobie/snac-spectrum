# snac-spectrum
Novel real time pitch tracking and harmonic decomposition based on a combined time/frequency domain approach

![Overview of system](/docs/snac-spectrum-overview.png)

# TODO
* Investigate more efficient real FFT algortihms
* Implement the flowchart logic and data flows
* Implement an algorithm for determining actual fundamental based on FFT

# Longer term ideas
* Pitch estimation immediately after attack (striking string) based on first max -> first zero crossing
* Investigate bundling an AI: Lookback after an attack which trains an AI to estimate next attack's pitch
* Investigate AI in general, once you have a stable and correct algo, it can be augmented by training an AI while in operation
* Fundamental determination based on AI taking FFT coefficients as input (f, f2, f3, ...) -> f/f2/f4 fundamental determination
* String cross contamination cleanup (mics are near each other and bound to bleed some)