// cmd/client/main.go
package main

import (
    "context"
    "log"
    "time"

    pb "github.com/Shridhar2104/logilo/shipment/proto/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)
func main() {
    conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewShipmentServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

// cmd/client/main.go
log.Println("\n=== Testing Calculate Rates ===")
rateResp, err := client.CalculateRates(ctx, &pb.RateRequest{
    OriginPincode:      "421503",
    DestinationPincode: "201002",
    Weight:             1,
    Length:             10,
    Width:              10,
    Height:             10,
    PaymentMode:        "prepaid",  // must be lowercase
    CollectableAmount:  10,
    CourierCodes:       []string{"XPRESSBEES"},
})
logResponse("Calculate Rates", rateResp, err)
    // // Test 2: Check Serviceability
    // log.Println("\nTesting Check Serviceability...")
    // availReq := &pb.AvailabilityRequest{
    //     OriginPincode:      "421503",
    //     DestinationPincode: "201002",
    //     Weight:             1,
    //     PaymentMode:        "prepaid",
    // }
    // availResp, err := client.GetAvailableCouriers(ctx, availReq)
    // logResponse("Serviceability", availResp, err)

    log.Println("\nTesting Create Shipment...")
    shipReq := &pb.CreateShipmentRequest{
        CourierCode:       "XPRESSBEES",    // Correct courier code
        OrderNumber:       "TEST123",
        PaymentType:       "prepaid",
        PackageWeight:     400,             // Weight in grams    
        PackageLength:     10,
        PackageBreadth:    10,
        PackageHeight:     10,
        OrderAmount:       1000.0,
        CollectableAmount: 0,               // 0 for prepaid
       // ShippingCharges:   177.0,           // From rate calculation
        
        Consignee: &pb.Address{
            Name:         "Test Customer",
            CompanyName:  "Test Company",
            AddressLine1: "123 Test Street",
            AddressLine2: "Near Test Location",
            City:         "Mumbai",
            State:        "Maharashtra",
            Pincode:      "400001",
            Phone:        "9999999999",
            Email:        "test@example.com",
        },
        
        Pickup: &pb.Address{
            Name:         "Test Contact",
            CompanyName:  "Test Warehouse",     // This is warehouse_name
            AddressLine1: "456 Sender Street",
            AddressLine2: "Near Sender Location",
            City:         "Delhi",
            State:        "Delhi",
            Pincode:      "110001",
            Phone:        "8888888888",
        },
        
        Items: []*pb.OrderItem{
            {
                Name:     "Test Product",
                Sku:      "SKU123",
                Quantity: 1,
                Price:    1000.0,
            },
        },
    }
    
    log.Printf("Creating shipment with request: %+v", shipReq)
    shipResp, err := client.CreateShipment(ctx, shipReq)
    if err != nil {
        log.Printf("Create shipment failed with error: %v", err)
        return
    }
    
    if !shipResp.Success {
        log.Printf("Create shipment failed: %s", shipResp.Error)
        return
    }
    
    log.Printf("Shipment created successfully. AWB: %s", shipResp.TrackingId)
    
    // Track shipment if AWB is present
    if shipResp.TrackingId != "" {
        log.Println("\nTesting Track Shipment...")
        trackReq := &pb.TrackingRequest{
            CourierCode: "XPRESSBEES",
            TrackingId:  shipResp.TrackingId,
        }
        trackResp, err := client.TrackShipment(ctx, trackReq)
        logResponse("Track Shipment", trackResp, err)
    }
}
func logResponse(operation string, response interface{}, err error) {
    if err != nil {
        log.Printf("%s failed: %v", operation, err)
    } else {
        log.Printf("%s Response: %+v", operation, response)
    }
}