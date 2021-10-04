/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.spark.models;

import java.util.Objects;
import java.util.Arrays;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;

/**
 * V1alpha1SparkClusterSpecImage
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecImage {
  public static final String SERIALIZED_NAME_PULL_POLICY = "pullPolicy";
  @SerializedName(SERIALIZED_NAME_PULL_POLICY)
  private String pullPolicy;

  public static final String SERIALIZED_NAME_REGISTRY = "registry";
  @SerializedName(SERIALIZED_NAME_REGISTRY)
  private String registry;

  public static final String SERIALIZED_NAME_REPOSITORY = "repository";
  @SerializedName(SERIALIZED_NAME_REPOSITORY)
  private String repository;

  public static final String SERIALIZED_NAME_TAG = "tag";
  @SerializedName(SERIALIZED_NAME_TAG)
  private String tag;


  public V1alpha1SparkClusterSpecImage pullPolicy(String pullPolicy) {
    
    this.pullPolicy = pullPolicy;
    return this;
  }

   /**
   * Get pullPolicy
   * @return pullPolicy
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getPullPolicy() {
    return pullPolicy;
  }


  public void setPullPolicy(String pullPolicy) {
    this.pullPolicy = pullPolicy;
  }


  public V1alpha1SparkClusterSpecImage registry(String registry) {
    
    this.registry = registry;
    return this;
  }

   /**
   * Get registry
   * @return registry
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getRegistry() {
    return registry;
  }


  public void setRegistry(String registry) {
    this.registry = registry;
  }


  public V1alpha1SparkClusterSpecImage repository(String repository) {
    
    this.repository = repository;
    return this;
  }

   /**
   * Get repository
   * @return repository
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getRepository() {
    return repository;
  }


  public void setRepository(String repository) {
    this.repository = repository;
  }


  public V1alpha1SparkClusterSpecImage tag(String tag) {
    
    this.tag = tag;
    return this;
  }

   /**
   * Get tag
   * @return tag
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getTag() {
    return tag;
  }


  public void setTag(String tag) {
    this.tag = tag;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecImage v1alpha1SparkClusterSpecImage = (V1alpha1SparkClusterSpecImage) o;
    return Objects.equals(this.pullPolicy, v1alpha1SparkClusterSpecImage.pullPolicy) &&
        Objects.equals(this.registry, v1alpha1SparkClusterSpecImage.registry) &&
        Objects.equals(this.repository, v1alpha1SparkClusterSpecImage.repository) &&
        Objects.equals(this.tag, v1alpha1SparkClusterSpecImage.tag);
  }

  @Override
  public int hashCode() {
    return Objects.hash(pullPolicy, registry, repository, tag);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecImage {\n");
    sb.append("    pullPolicy: ").append(toIndentedString(pullPolicy)).append("\n");
    sb.append("    registry: ").append(toIndentedString(registry)).append("\n");
    sb.append("    repository: ").append(toIndentedString(repository)).append("\n");
    sb.append("    tag: ").append(toIndentedString(tag)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

