// cmd/test/main.go
package main

import (
    "context"
    "log"
    "time"

    pb "github.com/Shridhar2104/logilo/shipment/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewShipmentServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

 // cmd/client/main.go
log.Println("\n=== Testing Calculate Rates ===")
rateResp, err := client.CalculateRates(ctx, &pb.RateRequest{
    OriginPincode:      "400001",
    DestinationPincode: "110001",
    Weight:             1,
    Length:             10,
    Width:              10,
    Height:             10,
    PaymentMode:        "prepaid",  // must be lowercase
    CollectableAmount:  0,
    CourierCodes:       []string{"XPRESSBEES"},
})
logResponse("Calculate Rates", rateResp, err)
    // Test 2: Check Serviceability
    log.Println("\n=== Testing Check Serviceability ===")
    availResp, err := client.GetAvailableCouriers(ctx, &pb.AvailabilityRequest{
        OriginPincode:      "400001",
        DestinationPincode: "110001",
        Weight:             1,
        PaymentMode:        "prepaid",
    })
    logResponse("Check Serviceability", availResp, err)

    // Test 3: Create Shipment
    log.Println("\n=== Testing Create Shipment ===")
    shipResp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
        CourierCode:    "XPRESSBEES",
        OrderNumber:    "TEST" + time.Now().Format("20060102150405"),
        PaymentType:    "prepaid",
        PackageWeight:  1,
        PackageLength:  10,
        PackageBreadth: 10,
        PackageHeight:  10,
        OrderAmount:    1000.0,
        CollectableAmount: 0,
        //ShippingCharges:   100,
        AutoPickup:        true,
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
            Name:         "Test Sender",
            CompanyName:  "Test Sender Company",
            AddressLine1: "456 Sender Street",
            AddressLine2: "Near Sender Location",
            City:         "Delhi",
            State:        "Delhi",
            Pincode:      "110001",
            Phone:        "8888888888",
            Email:        "sender@example.com",
            Gstin:        "GSTIN123456",
        },
        Items: []*pb.OrderItem{
            {
                Name:     "Test Product",
                Sku:      "SKU123",
                Quantity: 1,
                Price:    1000.0,
            },
        },
    })
    logResponse("Create Shipment", shipResp, err)

    // var trackingID string
    // if err == nil && shipResp.TrackingId != "" {
    //     trackingID = shipResp.TrackingId

    //     // Test 4: Track Shipment
    //     log.Println("\n=== Testing Track Shipment ===")
    //     trackResp, err := client.TrackShipment(ctx, &pb.TrackingRequest{
    //         CourierCode: "XPRESSBEES",
    //         TrackingId:  trackingID,
    //     })
    //     logResponse("Track Shipment", trackResp, err)

    //     // Test 5: Get NDR List
    //     log.Println("\n=== Testing Get NDR List ===")
    //     ndrResp, err := client.GetNDRList(ctx, &pb.NDRListRequest{
    //         CourierCode: "XPRESSBEES",
    //         Page:        1,
    //         Limit:       10,
    //     })
    //     logResponse("Get NDR List", ndrResp, err)

    //     // Test 6: Update NDR (if any NDR exists)
    //     if err == nil && len(ndrResp.Details) > 0 {
    //         log.Println("\n=== Testing Update NDR ===")
    //         updateResp, err := client.UpdateNDR(ctx, &pb.UpdateNDRRequest{
    //             CourierCode: "XPRESSBEES",
    //             Actions: []*pb.NDRAction{
    //                 {
    //                     AwbNumber: ndrResp.Details[0].AwbNumber,
    //                     Action:    "REATTEMPT",
    //                     Remarks:   "Please attempt delivery again",
    //                 },
    //             },
    //         })
    //         logResponse("Update NDR", updateResp, err)
    //     }

    //     // Test 7: Cancel Shipment
    //     log.Println("\n=== Testing Cancel Shipment ===")
    //     cancelResp, err := client.CancelShipment(ctx, &pb.CancelRequest{
    //         CourierCode: "XPRESSBEES",
    //         TrackingId:  trackingID,
    //     })
    //     logResponse("Cancel Shipment", cancelResp, err)
    // }
}

func logResponse(operation string, response interface{}, err error) {
    if err != nil {
        log.Printf("❌ %s failed: %v", operation, err)
    } else {
        log.Printf("✅ %s successful: %+v", operation, response)
    }
}