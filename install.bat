@ECHO OFF
for /f "tokens=1,* delims=:" %%A in ('curl -ks https://api.github.com/repos/KoperStudio/KoperManager/releases/latest ^| find "browser_download_url"') do (
    curl -kOL %%B
)
