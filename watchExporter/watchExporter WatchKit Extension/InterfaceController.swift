//
//  InterfaceController.swift
//  watchExporter WatchKit Extension
//
//  Created by Juda, Jacek on 3/13/20.
//  Copyright © 2020 Juda, Jacek. All rights reserved.
//

import WatchKit
import Foundation
import NIO
import GRPC
import SwiftProtobuf

class InterfaceController: WKInterfaceController {
    
    var client: Pb_CollectorClient!
    var msgNum: Int32 = 0
    
    override func awake(withContext context: Any?) {
        super.awake(withContext: context)
        
        // Configure interface objects here.
    }

    override func willActivate() {
        // This method is called when watch view controller is about to be visible to user
        super.willActivate()
    }
    
    override func didDeactivate() {
        // This method is called when watch view controller is no longer visible
        super.didDeactivate()
    }
    @IBOutlet weak var connect: WKInterfaceButton!
    @IBAction func Connect() {
        self.connect.setBackgroundColor(UIColor(red: 205, green: 144, blue: 55,alpha: 0.3 ))
        let group = MultiThreadedEventLoopGroup(numberOfThreads: 1)


        let channel = ClientConnection.insecure(group: group)
          .connect(host: "localhost", port: 15009)

        client = Pb_CollectorClient(channel: channel)
     
    }
    
    @IBOutlet weak var emmit: WKInterfaceButton!
    @IBAction func Emmit() {
        self.emmit.setBackgroundColor(UIColor(red: 205, green: 144, blue: 55,alpha: 0.3 ))
        emmitBio()
    }
    public func emmitBio() {
      print("→ Sending Bio Data")
      let options = CallOptions(timeout: .minutes(rounding: 1))
      let call = client.streamRecords(callOptions: options)
        
      call.response.whenSuccess { ack in
        print(
            "Finished trip with ack: \(ack.ack)."
        )
      }

      call.response.whenFailure { error in
        print("Emmitting Failed: \(error)")
      }

      call.status.whenComplete { _ in
        print("Finished Emmitting")
      }
    
      for x in 0..<self.msgNum {
        let bio : Pb_Biometrics = .with {
                $0.watchID = 1024967290
                $0.userName = "John Smith"
                $0.heartbeat = Int32.random(in: 1..<100)
                $0.temperature = Int32.random(in: 1..<100)
                $0.latitude = Float.random(in:-100..<100)
                $0.longitude = Float.random(in:-100..<100)
              }
        let item: Data
        do {
          item = try bio.serializedData()
        } catch {
         item = SwiftProtobuf.Internal.emptyData
        }
        let rc: Pb_Record = .with {
          $0.version = "v1"
          $0.topic = "apple_watch"
          $0.payload = item
        }
        call.sendMessage(rc, promise: nil)

        Thread.sleep(forTimeInterval: TimeInterval.random(in: 0.1..<0.5))
      }
      self.msgNum = 100
      call.sendEnd(promise: nil)

      _ = try! call.status.wait()
    }
}

