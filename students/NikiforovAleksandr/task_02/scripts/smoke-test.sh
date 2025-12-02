#!/bin/bash
echo "=== Smoke Test for Web16 Application ==="
echo "Testing namespace..."
kubectl get namespace app16

echo -e "\nTesting pods..."
kubectl get pods -n app16

echo -e "\nTesting services..."
kubectl get services -n app16

echo -e "\nTesting deployment..."
kubectl get deployment -n app16

echo -e "\nTesting ingress..."
kubectl get ingress -n app16

echo -e "\nChecking pod status..."
POD_NAME=$(kubectl get pods -n app16 -l app=web16 -o jsonpath='{.items[0].metadata.name}')
if [ -n "$POD_NAME" ]; then
    echo "First pod: $POD_NAME"
    kubectl logs -n app16 $POD_NAME --tail=5
    
    echo -e "\nTesting HTTP endpoint via port-forward..."
    kubectl port-forward -n app16 pod/$POD_NAME 8064:8064 &
    PF_PID=$!
    sleep 3
    
    curl -s http://localhost:8064/ | grep -o "ok.*" || echo "Health check failed"
    curl -s http://localhost:8064/health | grep -o "status.*" || echo "Health endpoint failed"
    
    kill $PF_PID
else
    echo "No pods found!"
fi

echo -e "\n=== Smoke test completed ==="