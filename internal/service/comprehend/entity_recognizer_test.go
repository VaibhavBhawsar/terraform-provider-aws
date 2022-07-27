package comprehend_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcomprehend "github.com/hashicorp/terraform-provider-aws/internal/service/comprehend"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccComprehendEntityRecognizer_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var entityrecognizer types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "data_access_role_arn", "aws_iam_role.test", "arn"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "comprehend", regexp.MustCompile(fmt.Sprintf(`entity-recognizer/%s$`, rName))),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.entity_types.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.annotations.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.augmented_manifests.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.data_format", string(types.EntityRecognizerDataFormatComprehendCsv)),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.documents.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "input_data_config.0.entity_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "language_code", "en"),
					resource.TestCheckResourceAttr(resourceName, "model_kms_key_id", ""),
					resource.TestCheckNoResourceAttr(resourceName, "model_policy"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "version_name", ""),
					resource.TestCheckResourceAttr(resourceName, "volume_kms_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "vpc_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var entityrecognizer types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					acctest.CheckResourceDisappears(acctest.Provider, tfcomprehend.ResourceEntityRecognizer(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_VersionName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var entityrecognizer types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	vName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	vName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_versionName(rName, vName1, "key", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "version_name", vName1),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "comprehend", regexp.MustCompile(fmt.Sprintf(`entity-recognizer/%s/version/%s$`, rName, vName1))),
					resource.TestCheckResourceAttrPair(resourceName, "data_access_role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckNoResourceAttr(resourceName, "model_policy"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEntityRecognizerConfig_versionName(rName, vName2, "key", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 2),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "version_name", vName2),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "comprehend", regexp.MustCompile(fmt.Sprintf(`entity-recognizer/%s/version/%s$`, rName, vName2))),
					resource.TestCheckResourceAttrPair(resourceName, "data_access_role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckNoResourceAttr(resourceName, "model_policy"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value2"),
				),
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_KMSKeys_CreateIDs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var entityrecognizer types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_kmsKeyIds(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttrPair(resourceName, "model_kms_key_id", "aws_kms_key.model", "key_id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_kms_key_id", "aws_kms_key.volume", "key_id"),
				),
			},
			{
				Config:   testAccEntityRecognizerConfig_kmsKeyARNs(rName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_KMSKeys_CreateARNs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var entityrecognizer types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_kmsKeyARNs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &entityrecognizer),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttrPair(resourceName, "model_kms_key_id", "aws_kms_key.model", "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_kms_key_id", "aws_kms_key.volume", "arn"),
				),
			},
			{
				Config:   testAccEntityRecognizerConfig_kmsKeyIds(rName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_KMSKeys_Update(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2, v3, v4 types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_kmsKeys_None(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v1),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "model_kms_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "volume_kms_key_id", ""),
				),
			},
			{
				Config: testAccEntityRecognizerConfig_kmsKeys_Set(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v2),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 2),
					resource.TestCheckResourceAttrPair(resourceName, "model_kms_key_id", "aws_kms_key.model", "key_id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_kms_key_id", "aws_kms_key.volume", "key_id"),
				),
			},
			{
				Config: testAccEntityRecognizerConfig_kmsKeys_Update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v3),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 3),
					resource.TestCheckResourceAttrPair(resourceName, "model_kms_key_id", "aws_kms_key.model2", "key_id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_kms_key_id", "aws_kms_key.volume2", "key_id"),
				),
			},
			{
				Config: testAccEntityRecognizerConfig_kmsKeys_None(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v4),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 4),
					resource.TestCheckResourceAttr(resourceName, "model_kms_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "volume_kms_key_id", ""),
				),
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_tags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2, v3 types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEntityRecognizerConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v1),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEntityRecognizerConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v2),
					testAccCheckEntityRecognizerNotRecreated(&v1, &v2),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccEntityRecognizerConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v3),
					testAccCheckEntityRecognizerNotRecreated(&v2, &v3),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccComprehendEntityRecognizer_DefaultTags_providerOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var v1, v2, v3 types.EntityRecognizerProperties
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_comprehend_entity_recognizer.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.ComprehendEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.ComprehendEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckEntityRecognizerDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("providerkey1", "providervalue1"),
					testAccEntityRecognizerConfig_tags0(rName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v1),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "providervalue1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags2("providerkey1", "providervalue1", "providerkey2", "providervalue2"),
					testAccEntityRecognizerConfig_tags0(rName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v2),
					testAccCheckEntityRecognizerNotRecreated(&v1, &v2),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "providervalue1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey2", "providervalue2"),
				),
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigDefaultTags_Tags1("providerkey1", "value1"),
					testAccEntityRecognizerConfig_tags0(rName),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEntityRecognizerExists(resourceName, &v3),
					testAccCheckEntityRecognizerNotRecreated(&v2, &v3),
					testAccCheckEntityRecognizerPublishedVersions(resourceName, 1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.providerkey1", "value1"),
				),
			},
		},
	})
}

// TODO: test deletion from in-error state. Try insufficient permissions to force error

// TODO: add test for catching, e.g. permission errors in training

func testAccCheckEntityRecognizerDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).ComprehendConn
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_comprehend_entity_recognizer" {
			continue
		}

		_, err := tfcomprehend.FindEntityRecognizerByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			if tfresource.NotFound(err) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Expected Comprehend Entity Recognizer to be destroyed, %s found", rs.Primary.ID)
	}

	return nil
}

func testAccCheckEntityRecognizerExists(name string, entityrecognizer *types.EntityRecognizerProperties) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Comprehend Entity Recognizer is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ComprehendConn
		ctx := context.Background()

		resp, err := tfcomprehend.FindEntityRecognizerByID(ctx, conn, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error describing Comprehend Entity Recognizer: %w", err)
		}

		*entityrecognizer = *resp

		return nil
	}
}

// func testAccCheckEntityRecognizerRecreated(before, after *types.EntityRecognizerProperties) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if entityRecognizerIdentity(before, after) {
// 			return fmt.Errorf("Comprehend Entity Recognizer not recreated")
// 		}

// 		return nil
// 	}
// }

func testAccCheckEntityRecognizerNotRecreated(before, after *types.EntityRecognizerProperties) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !entityRecognizerIdentity(before, after) {
			return fmt.Errorf("Comprehend Entity Recognizer recreated")
		}

		return nil
	}
}

func entityRecognizerIdentity(before, after *types.EntityRecognizerProperties) bool {
	return aws.ToTime(before.SubmitTime).Equal(aws.ToTime(after.SubmitTime))
}

func testAccCheckEntityRecognizerPublishedVersions(name string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Comprehend Entity Recognizer is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ComprehendConn
		ctx := context.Background()

		input := &comprehend.ListEntityRecognizersInput{
			Filter: &types.EntityRecognizerFilter{
				RecognizerName: aws.String(rs.Primary.Attributes["name"]),
			},
		}
		total := 0
		paginator := comprehend.NewListEntityRecognizersPaginator(conn, input)
		for paginator.HasMorePages() {
			output, err := paginator.NextPage(ctx)
			if err != nil {
				return err
			}
			total += len(output.EntityRecognizerPropertiesList)
		}
		return nil
	}
}

func testAccEntityRecognizerConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_versionName(rName, vName, key, value string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name         = %[1]q
  version_name = %[2]q

  data_access_role_arn = aws_iam_role.test.arn

  tags = {
    %[3]q = %[4]q
  }

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName, vName, key, value))
}

func testAccEntityRecognizerConfig_kmsKeyIds(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  model_kms_key_id  = aws_kms_key.model.key_id
  volume_kms_key_id = aws_kms_key.volume.key_id

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_iam_role_policy" "kms_keys" {
  role = aws_iam_role.test.name

  policy = data.aws_iam_policy_document.kms_keys.json
}

data "aws_iam_policy_document" "kms_keys" {
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.model.arn,
    ]
  }
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.volume.arn,
    ]
  }
}

resource "aws_kms_key" "model" {
  deletion_window_in_days = 7
}

resource "aws_kms_key" "volume" {
  deletion_window_in_days = 7
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_kmsKeyARNs(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  model_kms_key_id  = aws_kms_key.model.arn
  volume_kms_key_id = aws_kms_key.volume.arn

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_iam_role_policy" "kms_keys" {
  role = aws_iam_role.test.name

  policy = data.aws_iam_policy_document.kms_keys.json
}

data "aws_iam_policy_document" "kms_keys" {
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.model.arn,
    ]
  }
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.volume.arn,
    ]
  }
}

resource "aws_kms_key" "model" {
  deletion_window_in_days = 7
}

resource "aws_kms_key" "volume" {
  deletion_window_in_days = 7
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_kmsKeys_None(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_kmsKeys_Set(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  model_kms_key_id  = aws_kms_key.model.key_id
  volume_kms_key_id = aws_kms_key.volume.key_id

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_iam_role_policy" "kms_keys" {
  role = aws_iam_role.test.name

  policy = data.aws_iam_policy_document.kms_keys.json
}

data "aws_iam_policy_document" "kms_keys" {
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.model.arn,
    ]
  }
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.volume.arn,
    ]
  }
}

resource "aws_kms_key" "model" {
  deletion_window_in_days = 7
}

resource "aws_kms_key" "volume" {
  deletion_window_in_days = 7
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_kmsKeys_Update(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  model_kms_key_id  = aws_kms_key.model2.key_id
  volume_kms_key_id = aws_kms_key.volume2.key_id

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_iam_role_policy" "kms_keys" {
  role = aws_iam_role.test.name

  policy = data.aws_iam_policy_document.kms_keys.json
}

data "aws_iam_policy_document" "kms_keys" {
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.model2.arn,
    ]
  }
  statement {
    actions = [
      "*",
    ]

    resources = [
      aws_kms_key.volume2.arn,
    ]
  }
}

resource "aws_kms_key" "model2" {
  deletion_window_in_days = 7
}

resource "aws_kms_key" "volume2" {
  deletion_window_in_days = 7
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_tags0(rName string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  tags = {}

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName))
}

func testAccEntityRecognizerConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
  }

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName, tagKey1, tagValue1))
}

func testAccEntityRecognizerConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccEntityRecognizerBasicRoleConfig(rName),
		testAccEntityRecognizerS3BucketConfig(rName),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_comprehend_entity_recognizer" "test" {
  name = %[1]q

  data_access_role_arn = aws_iam_role.test.arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  language_code = "en"
  input_data_config {
    entity_types {
      type = "ENGINEER"
    }
    entity_types {
      type = "MANAGER"
    }

    documents {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.documents.id}"
    }

    entity_list {
      s3_uri = "s3://${aws_s3_bucket.test.bucket}/${aws_s3_object.entities.id}"
    }
  }

  depends_on = [
    aws_iam_role_policy.test
  ]
}

resource "aws_s3_object" "documents" {
  bucket = aws_s3_bucket.test.bucket
  key    = "documents.txt"
  source = "test-fixtures/entity_recognizer/documents.txt"
}

resource "aws_s3_object" "entities" {
  bucket = aws_s3_bucket.test.bucket
  key    = "entitylist.csv"
  source = "test-fixtures/entity_recognizer/entitylist.csv"
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccEntityRecognizerS3BucketConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.bucket

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.bucket

  rule {
    object_ownership = "BucketOwnerEnforced"
  }
}
`, rName)
}

func testAccEntityRecognizerBasicRoleConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "comprehend.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "test" {
  role = aws_iam_role.test.name

  policy = data.aws_iam_policy_document.role.json
}

data "aws_iam_policy_document" "role" {
  statement {
    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_bucket.test.arn}/*",
    ]
  }
  statement {
    actions = [
      "s3:ListBucket",
    ]

    resources = [
      aws_s3_bucket.test.arn,
    ]
  }
}
`, rName)
}
