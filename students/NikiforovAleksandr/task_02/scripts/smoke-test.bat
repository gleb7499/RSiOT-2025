@echo off
echo === Smoke Test for Web16 Application ===
echo Testing namespace...
kubectl get namespace app16

echo.
echo Testing pods...
kubectl get pods -n app16

echo.
echo Testing services...
kubectl get services -n app16

echo.
echo Testing deployment...
kubectl get deployment -n app16

echo.
echo Testing ingress...
kubectl get ingress -n app16

echo.
echo Checking pod status...
for /f "tokens=1" %%i in ('kubectl get pods -n app16 -l app^=web16 -o jsonpath^="{.items[0].metadata.name}" 2^>nul') do set POD_NAME=%%i

if "%POD_NAME%"=="" (
    echo No pods found!
) else (
    echo First pod: %POD_NAME%
    kubectl logs -n app16 %POD_NAME% --tail=5
    
    echo.
    echo Testing HTTP endpoint via port-forward...
    start /B kubectl port-forward -n app16 pod/%POD_NAME% 8064:8064
    timeout /t 3
    
    curl http://localhost:8064/
    curl http://localhost:8064/health
    
    taskkill /F /IM kubectl.exe >nul 2>&1
)

echo.
echo === Smoke test completed ===
pause