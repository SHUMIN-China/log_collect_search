package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
)


func Consume(consumer sarama.Consumer, topic string)  {
	partitionList, err := consumer.Partitions(topic)		//根据topic找到所有的分区
	if err != nil {
		fmt.Println(err)
		return
	}
	//遍历所有的分区
	wg := sync.WaitGroup{}
	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			fmt.Println(err)
		}
		defer pc.Close()
		//异步的消费数据
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			for msg := range partitionConsumer.Messages(){
				fmt.Printf("topic:%s,  partiotion:%v, offset:%v, value:%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}

func InitConsumer(conf *Config) sarama.Consumer {
	config := sarama.NewConfig()
	fmt.Println("====load kafka config success====")
	address := conf.Ip + ":" + conf.Port
	consumer, err := sarama.NewConsumer([]string{address}, config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return consumer
}