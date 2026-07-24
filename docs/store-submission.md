# Microsoft Store Submission — Plan and Windows Checklist

Decisions made 2026-07-23. This is the to-do list for the Windows machine.

## Decisions

- SimplyAuto goes on the Microsoft Store as a **free** app.
- We ship an **MSIX package**, built by running our WiX MSI through
  Microsoft's MSIX Packaging Tool. Microsoft signs the package, so we do
  not buy a code-signing certificate, and the Store handles updates.
- **Not UWP.** Go and Fyne cannot target it, and the UWP sandbox forbids
  the two things SimplyAuto exists to do: system-wide input hooks and
  sending input to other applications. A packaged desktop app gets all
  the Store benefits anyway.
- The privacy statement is live on the SimplyAuto website at
  `/privacy.html`. That URL goes into Partner Center.
- License info will be shown in a Help area inside the app. SimplyAuto
  itself is MIT.

## Step 1 — Get master building again

- Push the commits on this machine that fix the imports after the
  Refactor commit (`pkg/events` → `internal/events`, `assets` →
  `internal/assets`). As of d538606, the repo as pushed does not compile.
- Confirm with a clean clone or `git stash`: `go build ./...`
- While in there: delete the stray `internal/app/app..go` (double-dot
  filename, duplicate `App` type) and one of the duplicate
  `TODO.md` / `todo.md` files.

## Step 2 — Third-party license notices

Everything below gets bundled into one notices file and shown in the
app's Help area alongside our own MIT license.

Generate the file from the vendor directory (git bash):

```sh
for f in $(find vendor -maxdepth 4 -iname 'LICENSE*' -o -maxdepth 4 -iname 'COPYING*'); do
  echo "================================================================"
  echo "${f#vendor/}"
  echo "================================================================"
  cat "$f"
  echo
done > THIRD-PARTY-NOTICES.txt
```

The main ones for reference:

| Dependency | License |
|---|---|
| fyne.io/fyne/v2 | BSD-3-Clause |
| fyne.io/systray | Apache-2.0 |
| github.com/moutend/go-hook | MIT |
| golang.design/x/hotkey | MIT |
| golang.org/x/* | BSD-3-Clause |

Then add the Help tab (or an About dialog): app version, MIT license
text for SimplyAuto, and the notices file (embed it with the same
approach as `internal/assets`).

## Step 3 — Build the release exe

```
make build-debug        # sanity check first
make build VERSION=x.x.x
```

## Step 4 — Build the MSI with WiX

Keep the installer boring — the MSIX conversion just watches what it
does:

- Install `simplyauto.exe` to Program Files and add a Start Menu
  shortcut. Nothing else — no services, no drivers, no startup entries.
- Pick an UpgradeCode once and never change it.
- It must install silently (`msiexec /i simplyauto.msi /qn`) with no
  prompts — a plain WiX MSI does this by default; just don't add UI
  custom actions.

## Step 5 — Convert to MSIX

1. Reserve the app name first (Step 6) — the conversion needs identity
   values from Partner Center: package identity name, publisher ID
   (`CN=...`), and publisher display name. They're under
   Product management → Product identity.
2. Install the **MSIX Packaging Tool** (free, from the Store).
   Best run on a clean VM or fresh snapshot so it doesn't capture stray
   system noise.
3. Point it at the MSI, enter the identity values, let it convert. It
   declares the `runFullTrust` capability automatically — that's
   required and expected for a converted desktop app.
   The conversion is scripted: run elevated
   `MsixPackagingTool.exe create-package /template installer\msix-template.xml`
   (update the version in the template per release). Known quirk: the
   converted manifest picks up a spurious `Microsoft.WindowsAppRuntime`
   PackageDependency from the host machine; it's harmless — the Store
   auto-installs it — and stripping it would require makeappx from the
   Windows SDK.
4. Test the resulting `.msix` locally (Settings → For developers →
   Developer Mode lets you install unsigned test packages). Verify:
   - hotkeys F6/F9/F10/F11 work
   - recording captures keyboard and mouse; playback works
   - settings survive an app restart (registry writes are virtualized
     under MSIX — this is the thing most likely to behave differently)
   - saving and loading `.simplyauto` files works
   - uninstall from Start Menu is clean

## Step 6 — Partner Center

1. Register a developer account (individual, $19 one-time) at
   partner.microsoft.com.
2. Reserve the name **SimplyAuto**.
3. Create the submission:
   - Pricing: Free. Category: Utilities & tools.
   - Privacy policy URL: the live `/privacy.html` page.
   - Age rating: fill in the questionnaire (utility, no user content —
     comes out all-ages).
   - Upload the `.msix`. For the `runFullTrust` justification, one
     sentence: "Desktop automation utility; requires full trust to
     install input hooks for macro recording and to send synthetic
     input for auto-clicking."
   - Listing: description, logo, and the screenshots from the website
     repo (`images/screen1-3.png`).
4. Submit. Certification usually takes 1–3 business days.

## Listing wording — the one real rejection risk

Store policy bans products that enable cheating in online games. Auto
clickers are fine as general tools, so the listing must read that way:

- Describe it as a general-purpose automation tool: repetitive clicking
  tasks, testing, accessibility.
- Do not mention games, botting, AFK farming, or any game by name in
  the description, keywords, or screenshots.
- Being up front about what it does (records input locally, simulates
  clicks) helps with both certification and antivirus false positives —
  same framing already used on the website FAQ.
